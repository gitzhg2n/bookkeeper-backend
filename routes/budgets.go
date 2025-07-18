package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"bookkeeper-backend-go/models"
	"bookkeeper-backend-go/middleware"
)

func RegisterBudgetRoutes(r *mux.Router) {
	sub := r.PathPrefix("/budgets").Subrouter()
	sub.HandleFunc("/", createBudget).Methods("POST")
	sub.HandleFunc("/", listBudgets).Methods("GET")
	sub.HandleFunc("/{id}", getBudget).Methods("GET")
	sub.HandleFunc("/{id}", updateBudget).Methods("PUT")
	sub.HandleFunc("/{id}", deleteBudget).Methods("DELETE")
}

func createBudget(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	var budget models.Budget
	if err := json.NewDecoder(r.Body).Decode(&budget); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	if !middleware.CheckHouseholdOwnership(r.Context(), userID, budget.HouseholdID) {
		http.Error(w, "Forbidden: Not your household", http.StatusForbidden)
		return
	}
	if err := models.DB.Create(&budget).Error; err != nil {
		http.Error(w, "Error creating budget", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(budget)
}

func listBudgets(w http.ResponseWriter, r *http.Request) {
	householdIDs := r.Context().Value("householdIDs").([]uint)
	var budgets []models.Budget
	if err := models.DB.Where("household_id IN (?)", householdIDs).Find(&budgets).Error; err != nil {
		http.Error(w, "Error fetching budgets", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(budgets)
}

func getBudget(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var budget models.Budget
	if err := models.DB.First(&budget, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if !middleware.CheckHouseholdOwnership(r.Context(), userID, budget.HouseholdID) {
		http.Error(w, "Forbidden: Not your budget", http.StatusForbidden)
		return
	}
	json.NewEncoder(w).Encode(budget)
}

func updateBudget(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var budget models.Budget
	if err := models.DB.First(&budget, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if !middleware.CheckHouseholdOwnership(r.Context(), userID, budget.HouseholdID) {
		http.Error(w, "Forbidden: Not your budget", http.StatusForbidden)
		return
	}
	var payload models.Budget
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	budget.Name = payload.Name
	budget.Period = payload.Period
	budget.Category = payload.Category
	if err := models.DB.Save(&budget).Error; err != nil {
		http.Error(w, "Error updating budget", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(budget)
}

func deleteBudget(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var budget models.Budget
	if err := models.DB.First(&budget, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if !middleware.CheckHouseholdOwnership(r.Context(), userID, budget.HouseholdID) {
		http.Error(w, "Forbidden: Not your budget", http.StatusForbidden)
		return
	}
	if err := models.DB.Delete(&budget).Error; err != nil {
		http.Error(w, "Error deleting budget", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Budget deleted"})
}