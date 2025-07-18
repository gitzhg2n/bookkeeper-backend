package routes

import (
	"encoding/json"
	"net/http"
	"bookkeeper-backend-go/models"
	"github.com/gorilla/mux"
)

func RegisterIncomeSourceRoutes(r *mux.Router) {
	sub := r.PathPrefix("/incomeSources").Subrouter()
	sub.HandleFunc("", getIncomeSources).Methods("GET")
	// Add other handlers here
}

func getIncomeSources(w http.ResponseWriter, r *http.Request) {
	var sources []models.IncomeSource
	models.DB.Find(&sources)

	type IncomeSourceResponse struct {
		ID         uint    `json:"id"`
		Name       string  `json:"name"`
		Type       string  `json:"type"`
		Amount     float64 `json:"amount"`
		Frequency  string  `json:"frequency"`
		Notes      string  `json:"notes"`
		HouseholdID uint   `json:"householdId"`
	}

	resp := make([]IncomeSourceResponse, len(sources))
	for i, src := range sources {
		resp[i] = IncomeSourceResponse{
			ID:          src.ID,
			Name:        src.Name,
			Type:        src.Type,
			Amount:      src.Amount,
			Frequency:   src.Frequency,
			Notes:       src.Notes,
			HouseholdID: src.HouseholdID,
		}
	}

	json.NewEncoder(w).Encode(resp)
}