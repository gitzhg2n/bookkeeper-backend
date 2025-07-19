package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"bookkeeper-backend/models"
	"github.com/gorilla/mux"
	"bookkeeper-backend/middleware"
)

type AccountRequest struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Institution string  `json:"institution"`
	Balance     float64 `json:"balance"`
}

func RegisterAccountRoutes(r *mux.Router) {
	sub := r.PathPrefix("/accounts").Subrouter()
	sub.HandleFunc("", getAccounts).Methods("GET")
	sub.HandleFunc("", createAccount).Methods("POST")
	sub.HandleFunc("/{id}", updateAccount).Methods("PUT")
	sub.HandleFunc("/{id}", deleteAccount).Methods("DELETE")
}

func getAccounts(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.GetUserContext(r.Context())
	if userCtx == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	var accounts []models.Account
	// Get accounts for user's households or accounts directly owned
	if err := models.DB.Where("household_id IN ?", userCtx.HouseholdIDs).Find(&accounts).Error; err != nil {
		http.Error(w, "Failed to get accounts", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.GetUserContext(r.Context())
	if userCtx == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	var req AccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// For simplicity, use the first household ID - in real app, this should be specified
	if len(userCtx.HouseholdIDs) == 0 {
		http.Error(w, "No household found for user", http.StatusBadRequest)
		return
	}
	
	account := models.Account{
		Name:        req.Name,
		Type:        req.Type,
		Institution: req.Institution,
		Balance:     req.Balance,
		HouseholdID: userCtx.HouseholdIDs[0],
	}
	
	if err := models.DB.Create(&account).Error; err != nil {
		http.Error(w, "Failed to create account", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
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