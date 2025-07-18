package routes

import (
	"encoding/json"
	"net/http"
	"bookkeeper-backend-go/models"
	"github.com/gorilla/mux"
)

func RegisterTransactionRoutes(r *mux.Router) {
	sub := r.PathPrefix("/transactions").Subrouter()
	sub.HandleFunc("", getTransactions).Methods("GET")
	// Add other handlers here
}

func getTransactions(w http.ResponseWriter, r *http.Request) {
	var txs []models.Transaction
	models.DB.Preload("Account").Find(&txs)

	type TransactionResponse struct {
		ID          uint    `json:"id"`
		Date        string  `json:"date"`
		Account     string  `json:"account"`
		Category    string  `json:"category"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}

	resp := make([]TransactionResponse, len(txs))
	for i, tx := range txs {
		resp[i] = TransactionResponse{
			ID:          tx.ID,
			Date:        tx.Date.Format("2006-01-02"),
			Account:     tx.Account.Name,
			Category:    tx.Category,
			Amount:      tx.Amount,
			Description: tx.Description,
		}
	}

	json.NewEncoder(w).Encode(resp)
}