package routes

import (
	"encoding/json"
	"net/http"
	"bookkeeper-backend-go/models"
	"github.com/gorilla/mux"
)

func RegisterAccountRoutes(r *mux.Router) {
	sub := r.PathPrefix("/accounts").Subrouter()
	sub.HandleFunc("", getAccounts).Methods("GET")
	// Add other handlers here
}

func getAccounts(w http.ResponseWriter, r *http.Request) {
	var accounts []models.Account
	models.DB.Find(&accounts)

	type AccountResponse struct {
		ID      uint    `json:"id"`
		Name    string  `json:"name"`
		Type    string  `json:"type"`
		Balance float64 `json:"balance"`
	}

	resp := make([]AccountResponse, len(accounts))
	for i, acc := range accounts {
		resp[i] = AccountResponse{
			ID:      acc.ID,
			Name:    acc.Name,
			Type:    acc.Type,
			Balance: acc.Balance,
		}
	}

	json.NewEncoder(w).Encode(resp)
}