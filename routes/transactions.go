package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"bookkeeper-backend/internal/models"
	"bookkeeper-backend/middleware"

	"gorm.io/gorm"
)

type TransactionHandler struct {
	db *gorm.DB
	Notifications *db.NotificationStore
}

func NewTransactionHandler(db *gorm.DB, notifications *db.NotificationStore) *TransactionHandler {
	return &TransactionHandler{db: db, Notifications: notifications}
}

type createTransactionRequest struct {
	AmountCents int64   `json:"amount_cents"`
	Currency    string  `json:"currency"`
	CategoryID  *uint   `json:"category_id"`
	Memo        string  `json:"memo"`
	OccurredAt  *string `json:"occurred_at"`
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request, accountIDStr string) {
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		writeJSONError(r, w, "unauthorized", http.StatusUnauthorized)
		return
	}
	accID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		writeJSONError(r, w, "invalid account id", http.StatusBadRequest)
		return
	}
	var acc models.Account
	if err := h.db.First(&acc, uint(accID)).Error; err != nil {
		writeJSONError(r, w, "account not found", http.StatusNotFound)
		return
	}
	isMember, _ := userIsHouseholdMember(h.db, user.ID, acc.HouseholdID)
	if !isMember {
		writeJSONError(r, w, "forbidden", http.StatusForbidden)
		return
	}
	var req createTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(r, w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.AmountCents == 0 {
		writeJSONError(r, w, "amount_cents required", http.StatusBadRequest)
		return
	}
	if req.Currency == "" {
		req.Currency = acc.Currency
	}
	occ := time.Now()
	if req.OccurredAt != nil {
		if t, err := time.Parse(time.RFC3339, *req.OccurredAt); err == nil {
			occ = t
		}
	}
	trx := &models.Transaction{
		AccountID:   acc.ID,
		UserID:      &user.ID,
		AmountCents: req.AmountCents,
		Currency:    req.Currency,
		CategoryID:  req.CategoryID,
		Memo:        sanitizeString(req.Memo),
		OccurredAt:  occ,
	}
	if err := h.db.Create(trx).Error; err != nil {
		writeJSONError(r, w, "create failed", http.StatusInternalServerError)
		return
	}

	// Notification for large transaction
	threshold := int64(25000) // $250 in cents
	if user.Plan == "premium" || user.Plan == "selfhost" {
		settingsStore := db.UserSettingsStore{DB: h.db}
		if us, err := settingsStore.GetByUserID(user.ID); err == nil && us.LargeTransactionThreshold > 0 {
			threshold = us.LargeTransactionThreshold
		}
	}
	if req.AmountCents >= threshold {
		msg := "Large transaction detected: $" + fmt.Sprintf("%.2f", float64(req.AmountCents)/100)
		n := &models.Notification{
			UserID:  int64(user.ID),
			Type:    models.NotificationTypeTransaction,
			Message: msg,
			Read:    false,
			CreatedAt: time.Now(),
		}
		if h.Notifications != nil {
			h.Notifications.CreateNotification(r.Context(), n)
		}
	}

	// Goal progress notifications (50% and 100%)
	goalStore := db.GoalStore{DB: h.db}
	goals, _ := goalStore.ListByUser(user.ID)
	for _, goal := range goals {
		if goal.TargetCents == 0 {
			continue
		}
		progress := float64(goal.CurrentCents) / float64(goal.TargetCents)
		if progress >= 0.5 && progress < 1.0 {
			msg := "Goal '" + goal.Name + "' is 50% complete!"
			n := &models.Notification{
				UserID:  int64(user.ID),
				Type:    models.NotificationTypeGoal,
				Message: msg,
				Read:    false,
				CreatedAt: time.Now(),
			}
			if h.Notifications != nil {
				h.Notifications.CreateNotification(r.Context(), n)
			}
		} else if progress >= 1.0 {
			msg := "Goal '" + goal.Name + "' is complete!"
			n := &models.Notification{
				UserID:  int64(user.ID),
				Type:    models.NotificationTypeGoal,
				Message: msg,
				Read:    false,
				CreatedAt: time.Now(),
			}
			if h.Notifications != nil {
				h.Notifications.CreateNotification(r.Context(), n)
			}
		}
	}

	writeJSONSuccess(r, w, "created", trx)
}

func (h *TransactionHandler) List(w http.ResponseWriter, r *http.Request, accountIDStr string) {
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		writeJSONError(r, w, "unauthorized", http.StatusUnauthorized)
		return
	}
	accID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		writeJSONError(r, w, "invalid account id", http.StatusBadRequest)
		return
	}
	var acc models.Account
	if err := h.db.First(&acc, uint(accID)).Error; err != nil {
		writeJSONError(r, w, "account not found", http.StatusNotFound)
		return
	}
	isMember, _ := userIsHouseholdMember(h.db, user.ID, acc.HouseholdID)
	if !isMember {
		writeJSONError(r, w, "forbidden", http.StatusForbidden)
		return
	}

	query := h.db.Where("account_id = ?", acc.ID)

	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	if start != "" {
		if t, err := time.Parse(time.RFC3339, start); err == nil {
			query = query.Where("occurred_at >= ?", t)
		}
	}
	if end != "" {
		if t, err := time.Parse(time.RFC3339, end); err == nil {
			query = query.Where("occurred_at <= ?", t)
		}
	}

	var txs []models.Transaction
	query.Order("occurred_at desc").Limit(500).Find(&txs)
	writeJSONSuccess(r, w, "ok", txs)
}