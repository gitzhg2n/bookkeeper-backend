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
	sub.HandleFunc("", createGoal).Methods("POST")
}

type GoalRequest struct {
	Name     string  `json:"name"`
	Target   float64 `json:"target"`
	Progress float64 `json:"progress"`
}

func createGoal(w http.ResponseWriter, r *http.Request) {
	var req GoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid JSON"})
		return
	}
	if req.Name == "" || req.Target <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Name and positive target required"})
		return
	}
	if req.Progress < 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Progress cannot be negative"})
		return
	}
	goal := models.Goal{
		Name:     req.Name,
		Target:   req.Target,
		Progress: req.Progress,
	}
	if err := models.DB.Create(&goal).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to save goal"})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       goal.ID,
		"name":     goal.Name,
		"target":   goal.Target,
		"progress": goal.Progress,
	})
}

// ... existing getGoals code remains unchanged