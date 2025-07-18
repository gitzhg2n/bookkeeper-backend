package routes

import (
	"encoding/json"
	"net/http"
	"bookkeeper-backend-go/models"
	"github.com/gorilla/mux"
)

func RegisterBudgetRoutes(r *mux.Router) {
	sub := r.PathPrefix("/budgets").Subrouter()
	sub.HandleFunc("", getBudgets).Methods("GET")
	sub.HandleFunc("", createBudget).Methods("POST")
}

type BudgetRequest struct {
	Category string `json:"category"`
	Amount   int64  `json:"amount"`
}

func createBudget(w http.ResponseWriter, r *http.Request) {
	var req BudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid JSON"})
		return
	}
	if req.Category == "" || req.Amount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Category and positive amount required"})
		return
	}
	b := models.Budget{
		Category: req.Category,
		Amount:   req.Amount,
	}
	if err := models.DB.Create(&b).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to save budget"})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       b.ID,
		"category": b.Category,
		"amount":   b.Amount,
	})
}

// ... existing getBudgets code remains unchanged