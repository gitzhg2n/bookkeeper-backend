package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"bookkeeper-backend-go/models"
	"bookkeeper-backend-go/middleware"
)

func RegisterTransactionRoutes(r *mux.Router) {
	sub := r.PathPrefix("/transactions").Subrouter()
	sub.HandleFunc("/", createTransaction).Methods("POST")
	sub.HandleFunc("/", listTransactions).Methods("GET")
	sub.HandleFunc("/{id}", getTransaction).Methods("GET")
	sub.HandleFunc("/{id}", updateTransaction).Methods("PUT")
	sub.HandleFunc("/{id}", deleteTransaction).Methods("DELETE")
}

func createTransaction(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	var tx models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	if !middleware.CheckAccountOwnership(r.Context(), userID, tx.AccountID) {
		http.Error(w, "Forbidden: Not your account", http.StatusForbidden)
		return
	}
	if err := models.DB.Create(&tx).Error; err != nil {
		http.Error(w, "Error creating transaction", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tx)
}

func listTransactions(w http.ResponseWriter, r *http.Request) {
	accountIDs := r.Context().Value("accountIDs").([]uint)
	var transactions []models.Transaction
	if err := models.DB.Where("account_id IN (?)", accountIDs).Find(&transactions).Error; err != nil {
		http.Error(w, "Error fetching transactions", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(transactions)
}

func getTransaction(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var tx models.Transaction
	if err := models.DB.First(&tx, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if !middleware.CheckAccountOwnership(r.Context(), userID, tx.AccountID) {
		http.Error(w, "Forbidden: Not your transaction", http.StatusForbidden)
		return
	}
	json.NewEncoder(w).Encode(tx)
}

func updateTransaction(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var tx models.Transaction
	if err := models.DB.First(&tx, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if !middleware.CheckAccountOwnership(r.Context(), userID, tx.AccountID) {
		http.Error(w, "Forbidden: Not your transaction", http.StatusForbidden)
		return
	}
	var payload models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	tx.AccountID = payload.AccountID
	tx.Date = payload.Date
	tx.Category = payload.Category
	tx.Status = payload.Status
	tx.Amount = payload.Amount
	tx.Notes = payload.Notes
	if err := models.DB.Save(&tx).Error; err != nil {
		http.Error(w, "Error updating transaction", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tx)
}

func deleteTransaction(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var tx models.Transaction
	if err := models.DB.First(&tx, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if !middleware.CheckAccountOwnership(r.Context(), userID, tx.AccountID) {
		http.Error(w, "Forbidden: Not your transaction", http.StatusForbidden)
		return
	}
	if err := models.DB.Delete(&tx).Error; err != nil {
		http.Error(w, "Error deleting transaction", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction deleted"})
}