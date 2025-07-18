package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"bookkeeper-backend-go/models"
	"bookkeeper-backend-go/middleware"
)

func RegisterGoalRoutes(r *mux.Router) {
	sub := r.PathPrefix("/goals").Subrouter()
	sub.HandleFunc("/", createGoal).Methods("POST")
	sub.HandleFunc("/", listGoals).Methods("GET")
	sub.HandleFunc("/{id}", getGoal).Methods("GET")
	sub.HandleFunc("/{id}", updateGoal).Methods("PUT")
	sub.HandleFunc("/{id}", deleteGoal).Methods("DELETE")
}

func createGoal(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	var goal models.Goal
	if err := json.NewDecoder(r.Body).Decode(&goal); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	if !middleware.CheckHouseholdOwnership(r.Context(), userID, goal.HouseholdID) {
		http.Error(w, "Forbidden: Not your household", http.StatusForbidden)
		return
	}
	if err := models.DB.Create(&goal).Error; err != nil {
		http.Error(w, "Error creating goal", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(goal)
}

func listGoals(w http.ResponseWriter, r *http.Request) {
	householdIDs := r.Context().Value("householdIDs").([]uint)
	var goals []models.Goal
	if err := models.DB.Where("household_id IN (?)", householdIDs).Find(&goals).Error; err != nil {
		http.Error(w, "Error fetching goals", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(goals)
}

func getGoal(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var goal models.Goal
	if err := models.DB.First(&goal, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if !middleware.CheckHouseholdOwnership(r.Context(), userID, goal.HouseholdID) {
		http.Error(w, "Forbidden: Not your goal", http.StatusForbidden)
		return
	}
	json.NewEncoder(w).Encode(goal)
}

func updateGoal(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var goal models.Goal
	if err := models.DB.First(&goal, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if !middleware.CheckHouseholdOwnership(r.Context(), userID, goal.HouseholdID) {
		http.Error(w, "Forbidden: Not your goal", http.StatusForbidden)
		return
	}
	var payload models.Goal
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	goal.Name = payload.Name
	goal.TargetDate = payload.TargetDate
	goal.Category = payload.Category
	goal.Target = payload.Target
	goal.Progress = payload.Progress
	goal.Notes = payload.Notes
	if err := models.DB.Save(&goal).Error; err != nil {
		http.Error(w, "Error updating goal", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(goal)
}

func deleteGoal(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var goal models.Goal
	if err := models.DB.First(&goal, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if !middleware.CheckHouseholdOwnership(r.Context(), userID, goal.HouseholdID) {
		http.Error(w, "Forbidden: Not your goal", http.StatusForbidden)
		return
	}
	if err := models.DB.Delete(&goal).Error; err != nil {
		http.Error(w, "Error deleting goal", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Goal deleted"})
}