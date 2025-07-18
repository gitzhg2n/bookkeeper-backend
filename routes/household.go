package routes

import (
	"encoding/json"
	"net/http"
	"bookkeeper-backend-go/models"
	"github.com/gorilla/mux"
)

func RegisterHouseholdRoutes(r *mux.Router) {
	sub := r.PathPrefix("/household").Subrouter()
	sub.HandleFunc("/members", getHouseholdMembers).Methods("GET")
	// Add other handlers here
}

func getHouseholdMembers(w http.ResponseWriter, r *http.Request) {
	var members []models.HouseholdMember
	models.DB.Find(&members)

	type MemberResponse struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	resp := make([]MemberResponse, len(members))
	for i, m := range members {
		resp[i] = MemberResponse{
			ID:   m.ID,
			Name: m.Name,
		}
	}

	json.NewEncoder(w).Encode(resp)
}