package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"bookkeeper-backend-go/models"
	"github.com/gorilla/mux"
	"bookkeeper-backend-go/middleware"
)

func RegisterInvestmentRoutes(r *mux.Router) {
	sub := r.PathPrefix("/investments").Subrouter()
	sub.HandleFunc("", getInvestments).Methods("GET")
	sub.HandleFunc("", createInvestment).Methods("POST")
	sub.HandleFunc("/{id}", updateInvestment).Methods("PUT")
	sub.HandleFunc("/{id}", deleteInvestment).Methods("DELETE")
}

type InvestmentRequest struct {
	Name        string  `json:"name"`
	Value       float64 `json:"value"`
	Type        string  `json:"type"`
	Institution string  `json:"institution"`
}

func updateInvestment(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.GetUserContext(r.Context())
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if !middleware.CheckInvestmentOwnership(r.Context(), userCtx.ID, uint(id)) {
		http.Error(w, "Forbidden: Not your investment", http.StatusForbidden)
		return
	}
	var req InvestmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid JSON"})
		return
	}
	var inv models.Investment
	if err := models.DB.First(&inv, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Investment not found"})
		return
	}
	inv.Name = req.Name
	inv.Value = req.Value
	inv.Type = req.Type
	inv.Institution = req.Institution
	if err := models.DB.Save(&inv).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to update investment"})
		return
	}
	json.NewEncoder(w).Encode(inv)
}

func deleteInvestment(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.GetUserContext(r.Context())
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if !middleware.CheckInvestmentOwnership(r.Context(), userCtx.ID, uint(id)) {
		http.Error(w, "Forbidden: Not your investment", http.StatusForbidden)
		return
	}
	var inv models.Investment
	if err := models.DB.First(&inv, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Investment not found"})
		return
	}
	if err := models.DB.Delete(&inv).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to delete investment"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Investment deleted"})
}