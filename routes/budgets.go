package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"bookkeeper-backend/internal/models"
	"bookkeeper-backend/middleware"

	"gorm.io/gorm"
)

type BudgetHandler struct {
	db *gorm.DB
}

func NewBudgetHandler(db *gorm.DB) *BudgetHandler {
	return &BudgetHandler{db: db}
}

type createBudgetRequest struct {
	Month        string `json:"month"`          // YYYY-MM
	CategoryID   uint   `json:"category_id"`
	PlannedCents int64  `json:"planned_cents"`
}

func (h *BudgetHandler) Create(w http.ResponseWriter, r *http.Request, householdIDStr string) {
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		writeJSONError(r, w, "unauthorized", http.StatusUnauthorized)
		return
	}
	hID, valid := parseUintString(householdIDStr)
	if !valid {
		writeJSONError(r, w, "invalid household id", http.StatusBadRequest)
		return
	}
	isMember, _ := userIsHouseholdMember(h.db, user.ID, hID)
	if !isMember {
		writeJSONError(r, w, "forbidden", http.StatusForbidden)
		return
	}

	var req createBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	if !validMonth(req.Month) || req.CategoryID == 0 {
		writeJSONError(r, w, "invalid month or category_id", http.StatusBadRequest)
		return
	}
	b := &models.Budget{
		HouseholdID:  hID,
		Month:        req.Month,
		CategoryID:   req.CategoryID,
		PlannedCents: req.PlannedCents,
	}
	if err := h.db.Create(b).Error; err != nil {
		writeJSONError(r, w, "create failed (maybe duplicate)", http.StatusConflict)
		return
	}
	writeJSONSuccess(r, w, "created", b)
}

func (h *BudgetHandler) List(w http.ResponseWriter, r *http.Request, householdIDStr string) {
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		writeJSONError(r, w, "unauthorized", http.StatusUnauthorized)
		return
	}
	hID, valid := parseUintString(householdIDStr)
	if !valid {
		writeJSONError(r, w, "invalid household id", http.StatusBadRequest)
		return
	}
	isMember, _ := userIsHouseholdMember(h.db, user.ID, hID)
	if !isMember {
		writeJSONError(r, w, "forbidden", http.StatusForbidden)
		return
	}

	month := r.URL.Query().Get("month")
	if !validMonth(month) {
		writeJSONError(r, w, "month query param (YYYY-MM) required", http.StatusBadRequest)
		return
	}
	start, end := monthBounds(month)

	type row struct {
		ID           uint   `json:"id"`
		CategoryID   uint   `json:"category_id"`
		PlannedCents int64  `json:"planned_cents"`
		ActualCents  int64  `json:"actual_cents"`
		Month        string `json:"month"`
	}

	var budgets []models.Budget
	h.db.Where("household_id = ? AND month = ?", hID, month).Find(&budgets)

	results := make([]row, 0, len(budgets))
	for _, b := range budgets {
		var actual int64
		h.db.Model(&models.Transaction{}).
			Select("COALESCE(SUM(amount_cents),0)").
			Where("category_id = ? AND occurred_at >= ? AND occurred_at < ?", b.CategoryID, start, end).
			Scan(&actual)

		// Notification for budget threshold
		if actual >= b.PlannedCents && b.PlannedCents > 0 {
			msg := fmt.Sprintf("Budget reached for category %d: $%.2f spent of $%.2f", b.CategoryID, float64(actual)/100, float64(b.PlannedCents)/100)
			n := &models.Notification{
				UserID:  user.ID,
				Type:    models.NotificationTypeBudget,
				Message: msg,
				Read:    false,
				CreatedAt: time.Now(),
			}
			h.notifications.CreateNotification(r.Context(), n)
		}

		results = append(results, row{
			ID:           b.ID,
			CategoryID:   b.CategoryID,
			PlannedCents: b.PlannedCents,
			ActualCents:  actual,
			Month:        b.Month,
		})
	}

	writeJSONSuccess(r, w, "ok", results)
}

func validMonth(m string) bool {
	if len(m) != 7 {
		return false
	}
	_, err := time.Parse("2006-01", m)
	return err == nil
}

func monthBounds(m string) (time.Time, time.Time) {
	t, _ := time.Parse("2006-01", m)
	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	return start, end
}

func (h *BudgetHandler) Delete(w http.ResponseWriter, r *http.Request, householdIDStr, budgetIDStr string) {
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		writeJSONError(r, w, "unauthorized", http.StatusUnauthorized)
		return
	}
	hID, ok1 := parseUintString(householdIDStr)
	bID, ok2 := parseUintString(budgetIDStr)
	if !ok1 || !ok2 {
		writeJSONError(r, w, "invalid id", http.StatusBadRequest)
		return
	}
	isMember, _ := userIsHouseholdMember(h.db, user.ID, hID)
	if !isMember {
		writeJSONError(r, w, "forbidden", http.StatusForbidden)
		return
	}
	if err := h.db.Where("id = ? AND household_id = ?", bID, hID).Delete(&models.Budget{}).Error; err != nil {
		writeJSONError(r, w, "delete failed", http.StatusInternalServerError)
		return
	}
	writeJSONSuccess(r, w, "deleted", map[string]any{"id": bID})
}

func (h *BudgetHandler) Upsert(w http.ResponseWriter, r *http.Request, householdIDStr string) {
	// Optional convenience endpoint (PUT) to set planned amount
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		writeJSONError(r, w, "unauthorized", http.StatusUnauthorized)
		return
	}
	hID, valid := parseUintString(householdIDStr)
	if !valid {
		writeJSONError(r, w, "invalid household id", http.StatusBadRequest)
		return
	}
	isMember, _ := userIsHouseholdMember(h.db, user.ID, hID)
	if !isMember {
		writeJSONError(r, w, "forbidden", http.StatusForbidden)
		return
	}
	var req createBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	if !validMonth(req.Month) || req.CategoryID == 0 {
		writeJSONError(r, w, "invalid month or category", http.StatusBadRequest)
		return
	}
	var existing models.Budget
	err := h.db.Where("household_id = ? AND month = ? AND category_id = ?", hID, req.Month, req.CategoryID).
		First(&existing).Error
	if err == nil {
		// update
		existing.PlannedCents = req.PlannedCents
		if err := h.db.Save(&existing).Error; err != nil {
			writeJSONError(r, w, "update failed", http.StatusInternalServerError)
			return
		}
		writeJSONSuccess(r, w, "updated", existing)
		return
	}
	if err != gorm.ErrRecordNotFound {
		writeJSONError(r, w, "query error", http.StatusInternalServerError)
		return
	}
	// create
	b := &models.Budget{
		HouseholdID:  hID,
		Month:        req.Month,
		CategoryID:   req.CategoryID,
		PlannedCents: req.PlannedCents,
	}
	if err := h.db.Create(b).Error; err != nil {
		writeJSONError(r, w, "create failed", http.StatusInternalServerError)
		return
	}
	writeJSONSuccess(r, w, "created", b)
}

func (h *BudgetHandler) Summary(w http.ResponseWriter, r *http.Request, householdIDStr string) {
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		writeJSONError(r, w, "unauthorized", http.StatusUnauthorized)
		return
	}
	hID, valid := parseUintString(householdIDStr)
	if !valid {
		writeJSONError(r, w, "invalid household id", http.StatusBadRequest)
		return
	}
	isMember, _ := userIsHouseholdMember(h.db, user.ID, hID)
	if !isMember {
		writeJSONError(r, w, "forbidden", http.StatusForbidden)
		return
	}
	month := r.URL.Query().Get("month")
	if !validMonth(month) {
		writeJSONError(r, w, "month required YYYY-MM", http.StatusBadRequest)
		return
	}
	start, end := monthBounds(month)

	type row struct {
		CategoryID   uint   `json:"category_id"`
		PlannedCents int64  `json:"planned_cents"`
		ActualCents  int64  `json:"actual_cents"`
		Variance     int64  `json:"variance"`
		Month        string `json:"month"`
	}

	var budgets []models.Budget
	h.db.Where("household_id = ? AND month = ?", hID, month).Find(&budgets)

	out := make([]row, 0, len(budgets))
	for _, b := range budgets {
		var actual int64
		h.db.Model(&models.Transaction{}).
			Select("COALESCE(SUM(amount_cents),0)").
			Where("category_id = ? AND occurred_at >= ? AND occurred_at < ?", b.CategoryID, start, end).
			Scan(&actual)
		out = append(out, row{
			CategoryID:   b.CategoryID,
			PlannedCents: b.PlannedCents,
			ActualCents:  actual,
			Variance:     b.PlannedCents - actual,
			Month:        month,
		})
	}
	writeJSONSuccess(r, w, "ok", map[string]any{
		"month":   month,
		"items":   out,
		"summary": fmt.Sprintf("%d categories", len(out)),
	})
}