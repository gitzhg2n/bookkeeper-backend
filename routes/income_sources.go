package routes

import (
	"encoding/json"
	"net/http"
	"bookkeeper-backend/models"
	"github.com/gorilla/mux"
)

func RegisterIncomeSourceRoutes(r *mux.Router) {
	sub := r.PathPrefix("/incomeSources").Subrouter()
	sub.HandleFunc("", getIncomeSources).Methods("GET")
	sub.HandleFunc("", createIncomeSource).Methods("POST")
}

type IncomeSourceRequest struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Frequency   string  `json:"frequency"`
	Notes       string  `json:"notes"`
	HouseholdID uint    `json:"householdId"`
}

func getIncomeSources(w http.ResponseWriter, r *http.Request) {
	var incomeSources []models.IncomeSource
	if err := models.DB.Find(&incomeSources).Error; err != nil {
		http.Error(w, "Failed to get income sources", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(incomeSources)
}

func createIncomeSource(w http.ResponseWriter, r *http.Request) {
	var req IncomeSourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid JSON"})
		return
	}
	if req.Name == "" || req.Type == "" || req.Amount < 0 || req.Frequency == "" || req.HouseholdID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "All fields required; amount must be non-negative"})
		return
	}
	src := models.IncomeSource{
		Name:        req.Name,
		Type:        req.Type,
		Amount:      req.Amount,
		Frequency:   req.Frequency,
		Notes:       req.Notes,
		HouseholdID: req.HouseholdID,
	}
	if err := models.DB.Create(&src).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to save income source"})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":          src.ID,
		"name":        src.Name,
		"type":        src.Type,
		"amount":      src.Amount,
		"frequency":   src.Frequency,
		"notes":       src.Notes,
		"householdId": src.HouseholdID,
	})
}