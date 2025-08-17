package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bookkeeper-backend/internal/db"
	"bookkeeper-backend/internal/models"
	"bookkeeper-backend/middleware"
)

type InvestmentAlertHandler struct {
	Store *db.InvestmentAlertStore
}

// GET /investment-alerts - list all alerts for the user
func (h *InvestmentAlertHandler) List(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	alerts, err := h.Store.ListByUser(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "failed to fetch alerts", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(alerts)
}

// POST /investment-alerts - create a new alert
func (h *InvestmentAlertHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req models.InvestmentAlert
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	req.UserID = user.ID
	if err := h.Store.Create(r.Context(), &req); err != nil {
		http.Error(w, "failed to create alert", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// PUT /investment-alerts/{id} - update an alert
func (h *InvestmentAlertHandler) Update(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var req models.InvestmentAlert
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	req.ID = uint(id)
	req.UserID = user.ID
	if err := h.Store.Update(r.Context(), &req); err != nil {
		http.Error(w, "failed to update alert", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DELETE /investment-alerts/{id} - delete an alert
func (h *InvestmentAlertHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFrom(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.Store.Delete(r.Context(), user.ID, uint(id)); err != nil {
		http.Error(w, "failed to delete alert", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
