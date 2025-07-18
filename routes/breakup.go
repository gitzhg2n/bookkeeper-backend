package routes

import (
	"encoding/json"
	"net/http"
	"bookkeeper-backend-go/models"
	"github.com/gorilla/mux"
)

func RegisterBreakupRoutes(r *mux.Router) {
	sub := r.PathPrefix("/household/breakup").Subrouter()
	sub.HandleFunc("/", handleBreakup).Methods("POST")
}

func handleBreakup(w http.ResponseWriter, r *http.Request) {
	var req models.BreakupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	// TODO: Implement real breakup flow, updating household state, transactions, etc.
	json.NewEncoder(w).Encode(map[string]string{"message": "Breakup process initiated!"})
}