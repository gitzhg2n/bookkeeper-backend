package routes

import (
	"encoding/json"
	"net/http"
	"bookkeeper-backend-go/models"
	"github.com/gorilla/mux"
)

func RegisterInvestmentRoutes(r *mux.Router) {
	sub := r.PathPrefix("/investments").Subrouter()
	sub.HandleFunc("", getInvestments).Methods("GET")
	sub.HandleFunc("", createInvestment).Methods("POST")
}

type InvestmentRequest struct {
	Name   string  `json:"name"`
	Value  float64 `json:"value"`
	Type   string  `json:"type"`
	Change float64 `json:"change"`
}

func createInvestment(w http.ResponseWriter, r *http.Request) {
	var req InvestmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid JSON"})
		return
	}
	if req.Name == "" || req.Type == "" || req.Value < 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Name, type, and non-negative value required"})
		return
	}
	inv := models.Investment{
		Name:   req.Name,
		Value:  req.Value,
		Type:   req.Type,
		Change: req.Change,
	}
	if err := models.DB.Create(&inv).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to save investment"})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     inv.ID,
		"name":   inv.Name,
		"value":  inv.Value,
		"type":   inv.Type,
		"change": inv.Change,
	})
}

// ... existing getInvestments code remains unchanged