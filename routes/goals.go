package routes

import (
	"encoding/json"
	"net/http"
	"bookkeeper-backend-go/models"
	"github.com/gorilla/mux"
)

func RegisterGoalRoutes(r *mux.Router) {
	sub := r.PathPrefix("/goals").Subrouter()
	sub.HandleFunc("", getGoals).Methods("GET")
	// Add other handlers here
}

func getGoals(w http.ResponseWriter, r *http.Request) {
	var goals []models.Goal
	models.DB.Find(&goals)

	type GoalResponse struct {
		ID       uint    `json:"id"`
		Name     string  `json:"name"`
		Target   float64 `json:"target"`
		Progress float64 `json:"progress"`
	}

	resp := make([]GoalResponse, len(goals))
	for i, goal := range goals {
		resp[i] = GoalResponse{
			ID:       goal.ID,
			Name:     goal.Name,
			Target:   goal.Target,
			Progress: goal.Progress,
		}
	}

	json.NewEncoder(w).Encode(resp)
}