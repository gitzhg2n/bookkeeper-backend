package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bookkeeper-backend/internal/models"
	"bookkeeper-backend/middleware"

	"gorm.io/gorm"
)

type CategoryHandler struct {
	db *gorm.DB
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

type createCategoryRequest struct {
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request, householdIDStr string) {
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		writeJSONError(r, w, "unauthorized", http.StatusUnauthorized)
		return
	}
	hID, ok2 := parseUintString(householdIDStr)
	if !ok2 {
		writeJSONError(r, w, "invalid household id", http.StatusBadRequest)
		return
	}
	isMember, _ := userIsHouseholdMember(h.db, user.ID, hID)
	if !isMember {
		writeJSONError(r, w, "forbidden", http.StatusForbidden)
		return
	}
	var req createCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || sanitizeString(req.Name) == "" {
		writeJSONError(r, w, "invalid payload", http.StatusBadRequest)
		return
	}
	cat := &models.Category{
		HouseholdID: hID,
		Name:        sanitizeString(req.Name),
		ParentID:    req.ParentID,
	}
	if err := h.db.Create(cat).Error; err != nil {
		writeJSONError(r, w, "create failed", http.StatusInternalServerError)
		return
	}
	writeJSONSuccess(r, w, "created", cat)
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request, householdIDStr string) {
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		writeJSONError(r, w, "unauthorized", http.StatusUnauthorized)
		return
	}
	hID, ok2 := parseUintString(householdIDStr)
	if !ok2 {
		writeJSONError(r, w, "invalid household id", http.StatusBadRequest)
		return
	}
	isMember, _ := userIsHouseholdMember(h.db, user.ID, hID)
	if !isMember {
		writeJSONError(r, w, "forbidden", http.StatusForbidden)
		return
	}
	var cats []models.Category
	h.db.Where("household_id = ?", hID).Order("name asc").Find(&cats)
	writeJSONSuccess(r, w, "ok", cats)
}

// Simple stats endpoint: total categories for a household
func (h *CategoryHandler) Count(w http.ResponseWriter, r *http.Request, householdIDStr string) {
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		writeJSONError(r, w, "unauthorized", http.StatusUnauthorized)
		return
	}
	hID64, err := strconv.ParseUint(householdIDStr, 10, 32)
	if err != nil {
		writeJSONError(r, w, "invalid household id", http.StatusBadRequest)
		return
	}
	isMember, _ := userIsHouseholdMember(h.db, user.ID, uint(hID64))
	if !isMember {
		writeJSONError(r, w, "forbidden", http.StatusForbidden)
		return
	}
	var count int64
	h.db.Model(&models.Category{}).Where("household_id = ?", uint(hID64)).Count(&count)
	writeJSONSuccess(r, w, "ok", map[string]any{"count": count})
}