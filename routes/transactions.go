package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"bookkeeper-backend-go/models"
	"github.com/gorilla/mux"
	"bookkeeper-backend-go/middleware"
)

func RegisterTransactionRoutes(r *mux.Router) {
	sub := r.PathPrefix("/transactions").Subrouter()
	sub.HandleFunc("", getTransactions).Methods("GET")
	sub.HandleFunc("", createTransaction).Methods("POST")
	sub.HandleFunc("/{id}", updateTransaction).Methods("PUT")
	sub.HandleFunc("/{id}", deleteTransaction).Methods("DELETE")
}

type TransactionRequest struct {
	Date        string  `json:"date"`
	AccountID   uint    `json:"accountId"`
	Category    string  `json:"category"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

func updateTransaction(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.GetUserContext(r.Context())
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if !middleware.CheckTransactionOwnership(r.Context(), userCtx.ID, uint(id)) {
		http.Error(w, "Forbidden: Not your transaction", http.StatusForbidden)
		return
	}
	var req TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid JSON"})
		return
	}
	var tx models.Transaction
	if err := models.DB.First(&tx, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Transaction not found"})
		return
	}
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid date format"})
		return
	}
	tx.Date = parsedDate
	tx.AccountID = req.AccountID
	tx.Category = req.Category
	tx.Amount = req.Amount
	tx.Description = req.Description
	if err := models.DB.Save(&tx).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to update transaction"})
		return
	}
	json.NewEncoder(w).Encode(tx)
}

func deleteTransaction(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.GetUserContext(r.Context())
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if !middleware.CheckTransactionOwnership(r.Context(), userCtx.ID, uint(id)) {
		http.Error(w, "Forbidden: Not your transaction", http.StatusForbidden)
		return
	}
	var tx models.Transaction
	if err := models.DB.First(&tx, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Transaction not found"})
		return
	}
	if err := models.DB.Delete(&tx).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to delete transaction"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction deleted"})
}