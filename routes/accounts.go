package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"bookkeeper-backend-go/models"
	"github.com/gorilla/mux"
	"bookkeeper-backend-go/middleware"
)

// ... existing RegisterAccountRoutes, AccountRequest, createAccount, getAccounts ...

func RegisterAccountRoutes(r *mux.Router) {
	sub := r.PathPrefix("/accounts").Subrouter()
	sub.HandleFunc("", getAccounts).Methods("GET")
	sub.HandleFunc("", createAccount).Methods("POST")
	sub.HandleFunc("/{id}", updateAccount).Methods("PUT")
	sub.HandleFunc("/{id}", deleteAccount).Methods("DELETE")
}

func updateAccount(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.GetUserContext(r.Context())
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if !middleware.CheckAccountOwnership(r.Context(), userCtx.ID, uint(id)) {
		http.Error(w, "Forbidden: Not your account", http.StatusForbidden)
		return
	}

	var req AccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid JSON"})
		return
	}
	var acc models.Account
	if err := models.DB.First(&acc, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Account not found"})
		return
	}
	// Only update allowed fields
	acc.Name = req.Name
	acc.Type = req.Type
	acc.Institution = req.Institution
	acc.Balance = req.Balance
	if err := models.DB.Save(&acc).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to update account"})
		return
	}
	json.NewEncoder(w).Encode(acc)
}

func deleteAccount(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.GetUserContext(r.Context())
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if !middleware.CheckAccountOwnership(r.Context(), userCtx.ID, uint(id)) {
		http.Error(w, "Forbidden: Not your account", http.StatusForbidden)
		return
	}
	var acc models.Account
	if err := models.DB.First(&acc, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Account not found"})
		return
	}
	if err := models.DB.Delete(&acc).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to delete account"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Account deleted"})
}