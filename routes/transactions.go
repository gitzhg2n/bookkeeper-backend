package routes

import (
	"encoding/json"
	"net/http"
	"time"
	"bookkeeper-backend-go/models"
	"github.com/gorilla/mux"
)

func RegisterTransactionRoutes(r *mux.Router) {
	sub := r.PathPrefix("/transactions").Subrouter()
	sub.HandleFunc("", getTransactions).Methods("GET")
	sub.HandleFunc("", createTransaction).Methods("POST")
}

type TransactionRequest struct {
	Date        string  `json:"date"`
	AccountID   uint    `json:"accountId"`
	Category    string  `json:"category"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

func createTransaction(w http.ResponseWriter, r *http.Request) {
	var req TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid JSON"})
		return
	}
	if req.Date == "" || req.AccountID == 0 || req.Category == "" || req.Amount == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Date, accountId, category, and amount required"})
		return
	}
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid date format (YYYY-MM-DD required)"})
		return
	}
	tx := models.Transaction{
		Date:        parsedDate,
		AccountID:   req.AccountID,
		Category:    req.Category,
		Amount:      req.Amount,
		Description: req.Description,
	}
	if err := models.DB.Create(&tx).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to save transaction"})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":          tx.ID,
		"date":        tx.Date.Format("2006-01-02"),
		"accountId":   tx.AccountID,
		"category":    tx.Category,
		"amount":      tx.Amount,
		"description": tx.Description,
	})
}

// ... existing getTransactions code remains unchanged