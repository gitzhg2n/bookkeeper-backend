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
	Name     string `json:"name"`
	Category string `json:"category"`
	Amount   int64  `json:"amount"`
	Period   string `json:"period"`
}

func getBudgets(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserContext(r.Context())
	var budgets []models.Budget
	models.DB.Where("user_id = ?", user.ID).Find(&budgets)
	json.NewEncoder(w).Encode(budgets)
}

func createBudget(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserContext(r.Context())
	var req BudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Category == "" || req.Amount <= 0 || req.Period == "" {
		http.Error(w, "Missing/invalid fields", http.StatusBadRequest)
		return
	}
	budget := models.Budget{
		UserID:   user.ID,
		Name:     req.Name,
		Category: req.Category,
		Amount:   req.Amount,
		Period:   req.Period,
	}
	if err := models.DB.Create(&budget).Error; err != nil {
		http.Error(w, "Failed to create budget", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(budget)
}

func updateBudget(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserContext(r.Context())
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var budget models.Budget
	if err := models.DB.First(&budget, id).Error; err != nil {
		http.Error(w, "Budget not found", http.StatusNotFound)
		return
	}
	if budget.UserID != user.ID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	var req BudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Category == "" || req.Amount <= 0 || req.Period == "" {
		http.Error(w, "Missing/invalid fields", http.StatusBadRequest)
		return
	}
	budget.Name = req.Name
	budget.Category = req.Category
	budget.Amount = req.Amount
	budget.Period = req.Period
	if err := models.DB.Save(&budget).Error; err != nil {
		http.Error(w, "Failed to update budget", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(budget)
}

func deleteBudget(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserContext(r.Context())
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var budget models.Budget
	if err := models.DB.First(&budget, id).Error; err != nil {
		http.Error(w, "Budget not found", http.StatusNotFound)
		return
	}
	if budget.UserID != user.ID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if err := models.DB.Delete(&budget).Error; err != nil {
		http.Error(w, "Failed to delete budget", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Budget deleted"})
}