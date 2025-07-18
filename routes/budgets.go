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
	// Add other handlers here
}

func getBudgets(w http.ResponseWriter, r *http.Request) {
	var budgets []models.Budget
	models.DB.Find(&budgets)

	type BudgetResponse struct {
		ID       uint   `json:"id"`
		Category string `json:"category"`
		Amount   int64  `json:"amount"`
	}

	resp := make([]BudgetResponse, len(budgets))
	for i, b := range budgets {
		resp[i] = BudgetResponse{
			ID:       b.ID,
			Category: b.Category,
			Amount:   b.Amount,
		}
	}

	json.NewEncoder(w).Encode(resp)
}