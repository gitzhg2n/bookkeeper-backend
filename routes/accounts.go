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
	sub.HandleFunc("", createAccount).Methods("POST")
}

type AccountRequest struct {
	Name    string  `json:"name"`
	Type    string  `json:"type"`
	Balance float64 `json:"balance"`
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	var req AccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid JSON"})
		return
	}
	if req.Name == "" || req.Type == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Name and type are required"})
		return
	}
	// Optionally, validate balance is non-negative
	if req.Balance < 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Balance cannot be negative"})
		return
	}
	acc := models.Account{
		Name:    req.Name,
		Type:    req.Type,
		Balance: req.Balance,
	}
	if err := models.DB.Create(&acc).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to save account"})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      acc.ID,
		"name":    acc.Name,
		"type":    acc.Type,
		"balance": acc.Balance,
	})
}

// ... existing getAccounts code remains unchanged