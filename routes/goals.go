package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"bookkeeper-backend-go/models"
	"github.com/gorilla/mux"
	"bookkeeper-backend-go/middleware"
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

func updateGoal(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.GetUserContext(r.Context())
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if !middleware.CheckGoalOwnership(r.Context(), userCtx.ID, uint(id)) {
		http.Error(w, "Forbidden: Not your goal", http.StatusForbidden)
		return
	}
	var req GoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid JSON"})
		return
	}
	var goal models.Goal
	if err := models.DB.First(&goal, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Goal not found"})
		return
	}
	goal.Name = req.Name
	goal.Target = req.Target
	goal.Progress = req.Progress
	goal.Category = req.Category
	goal.Notes = req.Notes
	if err := models.DB.Save(&goal).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to update goal"})
		return
	}
	json.NewEncoder(w).Encode(goal)
}

func deleteGoal(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.GetUserContext(r.Context())
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if !middleware.CheckGoalOwnership(r.Context(), userCtx.ID, uint(id)) {
		http.Error(w, "Forbidden: Not your goal", http.StatusForbidden)
		return
	}
	var goal models.Goal
	if err := models.DB.First(&goal, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Goal not found"})
		return
	}
	if err := models.DB.Delete(&goal).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to delete goal"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Goal deleted"})
}