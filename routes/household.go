package routes

import (
	"encoding/json"
	"net/http"
	"bookkeeper-backend/models"
	"github.com/gorilla/mux"
)

func RegisterHouseholdMemberRoutes(r *mux.Router) {
	sub := r.PathPrefix("/household").Subrouter()
	sub.HandleFunc("/members", getHouseholdMembers).Methods("GET")
	sub.HandleFunc("/members", addHouseholdMember).Methods("POST")
}

type MemberRequest struct {
	Name string `json:"name"`
}

func getHouseholdMembers(w http.ResponseWriter, r *http.Request) {
	// Basic implementation - in real app, this would get members from database
	members := []models.HouseholdMember{}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

func addHouseholdMember(w http.ResponseWriter, r *http.Request) {
	var req MemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid JSON"})
		return
	}
	if req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Name required"})
		return
	}
	member := models.HouseholdMember{
		Name: req.Name,
	}
	if err := models.DB.Create(&member).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to add member"})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":   member.ID,
		"name": member.Name,
	})
}

// ... existing getHouseholdMembers code remains unchanged