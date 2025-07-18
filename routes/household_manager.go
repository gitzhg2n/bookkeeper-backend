package routes

import (
	"encoding/json"
	"net/http"
	"bookkeeper-backend-go/models"
	"github.com/gorilla/mux"
)

func RegisterHouseholdManagerRoutes(r *mux.Router) {
	sub := r.PathPrefix("/household/members").Subrouter()
	sub.HandleFunc("/", listMembers).Methods("GET")
	sub.HandleFunc("/", addMember).Methods("POST")
}

func listMembers(w http.ResponseWriter, r *http.Request) {
	householdId := r.Context().Value("householdId").(uint)
	var members []models.HouseholdMember
	// TODO: Replace with real DB query
	members = append(members, models.HouseholdMember{ID: 1, Name: "Alice", HouseholdID: householdId})
	members = append(members, models.HouseholdMember{ID: 2, Name: "Bob", HouseholdID: householdId})
	json.NewEncoder(w).Encode(members)
}

func addMember(w http.ResponseWriter, r *http.Request) {
	var payload models.HouseholdMember
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	// TODO: Replace with real DB insert logic
	payload.ID = 3
	json.NewEncoder(w).Encode(payload)
}