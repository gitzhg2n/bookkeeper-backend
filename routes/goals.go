package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"bookkeeper-backend/models"
	"github.com/gorilla/mux"
	"bookkeeper-backend/middleware"
)

func RegisterGoalRoutes(r *mux.Router) {
	sub := r.PathPrefix("/goals").Subrouter()
	sub.HandleFunc("", getGoals).Methods("GET")
	sub.HandleFunc("", createGoal).Methods("POST")
	sub.HandleFunc("/{id}", updateGoal).Methods("PUT")
	sub.HandleFunc("/{id}", deleteGoal).Methods("DELETE")
}

type GoalRequest struct {
	Name     string  `json:"name"`
	Target   float64 `json:"target"`
	Progress float64 `json:"progress"`
	Category string  `json:"category"`
	Notes    string  `json:"notes"`
}

func getGoals(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserContext(r.Context())
	var goals []models.Goal
	models.DB.Where("user_id = ?", user.ID).Find(&goals)
	json.NewEncoder(w).Encode(goals)
}

func createGoal(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserContext(r.Context())
	var req GoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Category == "" || req.Target <= 0 {
		http.Error(w, "Missing/invalid fields", http.StatusBadRequest)
		return
	}
	goal := models.Goal{
		UserID:   user.ID,
		Name:     req.Name,
		Target:   req.Target,
		Progress: req.Progress,
		Category: req.Category,
		Notes:    req.Notes,
	}
	if err := models.DB.Create(&goal).Error; err != nil {
		http.Error(w, "Failed to create goal", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(goal)
}

func updateGoal(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserContext(r.Context())
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var goal models.Goal
	if err := models.DB.First(&goal, id).Error; err != nil {
		http.Error(w, "Goal not found", http.StatusNotFound)
		return
	}
	if goal.UserID != user.ID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	var req GoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Category == "" || req.Target <= 0 {
		http.Error(w, "Missing/invalid fields", http.StatusBadRequest)
		return
	}
	goal.Name = req.Name
	goal.Target = req.Target
	goal.Progress = req.Progress
	goal.Category = req.Category
	goal.Notes = req.Notes
	if err := models.DB.Save(&goal).Error; err != nil {
		http.Error(w, "Failed to update goal", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(goal)
}

func deleteGoal(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserContext(r.Context())
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var goal models.Goal
	if err := models.DB.First(&goal, id).Error; err != nil {
		http.Error(w, "Goal not found", http.StatusNotFound)
		return
	}
	if goal.UserID != user.ID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if err := models.DB.Delete(&goal).Error; err != nil {
		http.Error(w, "Failed to delete goal", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Goal deleted"})
}