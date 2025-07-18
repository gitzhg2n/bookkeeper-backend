package routes

import (
	"encoding/json"
	"net/http"
	"bookkeeper-backend-go/models"
	"github.com/gorilla/mux"
)

func RegisterInvestmentRoutes(r *mux.Router) {
	sub := r.PathPrefix("/investments").Subrouter()
	sub.HandleFunc("", getInvestments).Methods("GET")
	// Add other handlers here
}

func getInvestments(w http.ResponseWriter, r *http.Request) {
	var investments []models.Investment
	models.DB.Find(&investments)

	type InvestmentResponse struct {
		ID     uint    `json:"id"`
		Name   string  `json:"name"`
		Value  float64 `json:"value"`
		Type   string  `json:"type"`
		Change float64 `json:"change"`
	}

	resp := make([]InvestmentResponse, len(investments))
	for i, inv := range investments {
		resp[i] = InvestmentResponse{
			ID:     inv.ID,
			Name:   inv.Name,
			Value:  inv.Value,
			Type:   inv.Type,
			Change: inv.Change,
		}
	}

	json.NewEncoder(w).Encode(resp)
}