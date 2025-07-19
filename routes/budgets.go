package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"bookkeeper-backend-go/models"
	"github.com/gorilla/mux"
	"bookkeeper-backend-go/middleware"
)

func RegisterBudgetRoutes(r *mux.Router) {
	sub := r.PathPrefix("/budgets").Subrouter()
	sub.HandleFunc("", getBudgets).Methods("GET")
	sub.HandleFunc("", createBudget).Methods("POST")
	sub.HandleFunc("/{id}", updateBudget).Methods("PUT")
	sub.HandleFunc("/{id}", deleteBudget).Methods("DELETE")
}

type BudgetRequest struct {
	Category string `json:"category"`
	Amount   int64  `json:"amount"`
	Name     string `json:"name"`
	Period   string `json:"period"`
}

func updateBudget(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.GetUserContext(r.Context())
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if !middleware.CheckBudgetOwnership(r.Context(), userCtx.ID, uint(id)) {
		http.Error(w, "Forbidden: Not your budget", http.StatusForbidden)
		return
	}
	var req BudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid JSON"})
		return
	}
	var budget models.Budget
	if err := models.DB.First(&budget, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Budget not found"})
		return
	}
	budget.Category = req.Category
	budget.Amount = req.Amount
	budget.Name = req.Name
	budget.Period = req.Period
	if err := models.DB.Save(&budget).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to update budget"})
		return
	}
	json.NewEncoder(w).Encode(budget)
}

func deleteBudget(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.GetUserContext(r.Context())
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if !middleware.CheckBudgetOwnership(r.Context(), userCtx.ID, uint(id)) {
		http.Error(w, "Forbidden: Not your budget", http.StatusForbidden)
		return
	}
	var budget models.Budget
	if err := models.DB.First(&budget, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Budget not found"})
		return
	}
	if err := models.DB.Delete(&budget).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to delete budget"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Budget deleted"})
}