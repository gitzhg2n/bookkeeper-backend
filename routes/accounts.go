package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"bookkeeper-backend/internal/models"
	"bookkeeper-backend/middleware"

	"gorm.io/gorm"
)

type AccountHandler struct {
	db *gorm.DB
}

func NewAccountHandler(db *gorm.DB) *AccountHandler {
	return &AccountHandler{db: db}
}

type createAccountRequest struct {
	Name                string `json:"name"`
	Type                string `json:"type"`
	Currency            string `json:"currency"`
	OpeningBalanceCents int64  `json:"opening_balance_cents"`
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request, householdIDStr string) {
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
	member, _ := userIsHouseholdMember(h.db, user.ID, hID)
	if !member {
		writeJSONError(r, w, "forbidden", http.StatusForbidden)
		return
	}
	var req createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || sanitizeString(req.Name) == "" {
		writeJSONError(r, w, "invalid payload", http.StatusBadRequest)
		return
	}
	if req.Type == "" {
		req.Type = "checking"
	}
	if req.Currency == "" {
		req.Currency = "USD"
	}
	acc := &models.Account{
		HouseholdID:         hID,
		Name:                sanitizeString(req.Name),
		Type:                req.Type,
		Currency:            req.Currency,
		OpeningBalanceCents: req.OpeningBalanceCents,
	}
	if err := h.db.Create(acc).Error; err != nil {
		writeJSONError(r, w, "create failed", http.StatusInternalServerError)
		return
	}
	writeJSONSuccess(r, w, "created", acc)
}

func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request, householdIDStr string) {
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
	member, _ := userIsHouseholdMember(h.db, user.ID, hID)
	if !member {
		writeJSONError(r, w, "forbidden", http.StatusForbidden)
		return
	}
	var accounts []models.Account
	h.db.Where("household_id = ?", hID).Find(&accounts)
	writeJSONSuccess(r, w, "ok", accounts)
}

func (h *AccountHandler) ensureOwnership(userID, accountID uint) (*models.Account, bool) {
	var acc models.Account
	if err := h.db.First(&acc, accountID).Error; err != nil {
		return nil, false
	}
	isMember, _ := userIsHouseholdMember(h.db, userID, acc.HouseholdID)
	if !isMember {
		return nil, false
	}
	return &acc, true
}

func (h *AccountHandler) archiveAccount(acc *models.Account) {
	now := time.Now()
	acc.ArchivedAt = &now
	h.db.Save(acc)
}