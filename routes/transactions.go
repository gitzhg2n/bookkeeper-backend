package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"bookkeeper-backend/models"
	"github.com/gorilla/mux"
	"bookkeeper-backend/middleware"
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

func getTransactions(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserContext(r.Context())
	var transactions []models.Transaction
	models.DB.Where("user_id = ?", user.ID).Find(&transactions)
	json.NewEncoder(w).Encode(transactions)
}

func createTransaction(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserContext(r.Context())
	var req TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil || req.Amount == 0 || req.Category == "" || req.AccountID == 0 {
		http.Error(w, "Missing/invalid fields", http.StatusBadRequest)
		return
	}
	tx := models.Transaction{
		UserID:      user.ID,
		Date:        parsedDate,
		AccountID:   req.AccountID,
		Category:    req.Category,
		Amount:      req.Amount,
		Description: req.Description,
	}
	if err := models.DB.Create(&tx).Error; err != nil {
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tx)
}

func updateTransaction(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserContext(r.Context())
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var tx models.Transaction
	if err := models.DB.First(&tx, id).Error; err != nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}
	if tx.UserID != user.ID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	var req TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil || req.Amount == 0 || req.Category == "" || req.AccountID == 0 {
		http.Error(w, "Missing/invalid fields", http.StatusBadRequest)
		return
	}
	tx.Date = parsedDate
	tx.AccountID = req.AccountID
	tx.Category = req.Category
	tx.Amount = req.Amount
	tx.Description = req.Description
	if err := models.DB.Save(&tx).Error; err != nil {
		http.Error(w, "Failed to update transaction", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tx)
}

func deleteTransaction(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserContext(r.Context())
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var tx models.Transaction
	if err := models.DB.First(&tx, id).Error; err != nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}
	if tx.UserID != user.ID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if err := models.DB.Delete(&tx).Error; err != nil {
		http.Error(w, "Failed to delete transaction", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction deleted"})
}