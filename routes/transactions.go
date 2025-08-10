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

type TransactionHandler struct {
	db *gorm.DB
}

func NewTransactionHandler(db *gorm.DB) *TransactionHandler {
	return &TransactionHandler{db: db}
}

type createTransactionRequest struct {
	AmountCents int64   `json:"amount_cents"`
	Currency    string  `json:"currency"`
	CategoryID  *uint   `json:"category_id"`
	Memo        string  `json:"memo"`
	OccurredAt  *string `json:"occurred_at"` // RFC3339 or date
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request, accountIDStr string) {
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		writeJSONError(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	accID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		writeJSONError(w, "invalid account id", http.StatusBadRequest)
		return
	}
	var acc models.Account
	if err := h.db.First(&acc, uint(accID)).Error; err != nil {
		writeJSONError(w, "account not found", http.StatusNotFound)
		return
	}
	isMember, _ := userIsHouseholdMember(h.db, user.ID, acc.HouseholdID)
	if !isMember {
		writeJSONError(w, "forbidden", http.StatusForbidden)
		return
	}
	var req createTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.AmountCents == 0 {
		writeJSONError(w, "amount_cents required", http.StatusBadRequest)
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
		writeJSONError(w, "create failed", http.StatusInternalServerError)
		return
	}
	writeJSONSuccess(w, "created", trx)
}

func (h *TransactionHandler) List(w http.ResponseWriter, r *http.Request, accountIDStr string) {
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		writeJSONError(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	accID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		writeJSONError(w, "invalid account id", http.StatusBadRequest)
		return
	}
	var acc models.Account
	if err := h.db.First(&acc, uint(accID)).Error; err != nil {
		writeJSONError(w, "account not found", http.StatusNotFound)
		return
	}
	isMember, _ := userIsHouseholdMember(h.db, user.ID, acc.HouseholdID)
	if !isMember {
		writeJSONError(w, "forbidden", http.StatusForbidden)
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
	writeJSONSuccess(w, "ok", txs)
}