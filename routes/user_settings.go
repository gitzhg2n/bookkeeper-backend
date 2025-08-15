package routes

import (
	"encoding/json"
	"net/http"

	"bookkeeper-backend/internal/db"
	"bookkeeper-backend/middleware"
)

type UserSettingsHandler struct {
	Store *db.UserSettingsStore
}

// GET /user/settings - get current user's settings
func (h *UserSettingsHandler) Get(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	us, err := h.Store.GetByUserID(user.ID)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(us)
}

// POST /user/settings - update current user's settings
func (h *UserSettingsHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req struct {
		LargeTransactionThreshold int64 `json:"large_transaction_threshold"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	store := db.UserSettingsStore{DB: h.Store.DB}
	if err := store.Upsert(user.ID, req.LargeTransactionThreshold); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
