package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"bookkeeper-backend/internal/models"
	"bookkeeper-backend/middleware"

	"gorm.io/gorm"
)

type HouseholdHandler struct {
	db *gorm.DB
}

func NewHouseholdHandler(db *gorm.DB) *HouseholdHandler {
	return &HouseholdHandler{db: db}
}

type createHouseholdRequest struct {
	Name string `json:"name"`
}

func (h *HouseholdHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		writeJSONError(r, w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req createHouseholdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || sanitizeString(req.Name) == "" {
		writeJSONError(r, w, "invalid name", http.StatusBadRequest)
		return
	}
	house := &models.Household{
		Name:      sanitizeString(req.Name),
		CreatedBy: user.ID,
	}
	if err := h.db.Create(house).Error; err != nil {
		writeJSONError(r, w, "create failed", http.StatusInternalServerError)
		return
	}
	member := &models.HouseholdMember{
		HouseholdID: house.ID,
		UserID:      user.ID,
		Role:        "owner",
		CreatedAt:   time.Now(),
	}
	_ = h.db.Create(member).Error

	writeJSONSuccess(r, w, "created", map[string]any{
		"id":   house.ID,
		"name": house.Name,
	})
}

func (h *HouseholdHandler) List(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		writeJSONError(r, w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var results []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
		Role string `json:"role"`
	}
	h.db.Table("households h").
		Select("h.id, h.name, hm.role").
		Joins("JOIN household_members hm ON hm.household_id = h.id").
		Where("hm.user_id = ?", user.ID).
		Scan(&results)
	writeJSONSuccess(r, w, "ok", results)
}

func userIsHouseholdMember(db *gorm.DB, userID uint, householdID uint) (bool, string) {
	var hm models.HouseholdMember
	if err := db.Where("user_id = ? AND household_id = ?", userID, householdID).First(&hm).Error; err != nil {
		return false, ""
	}
	return true, hm.Role
}

func parseUintString(s string) (uint, bool) {
	id64, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, false
	}
	return uint(id64), true
}