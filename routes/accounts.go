package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"bookkeeper-backend/middleware"
	"bookkeeper-backend/models"
	"github.com/gorilla/mux"
)

type AccountRequest struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Institution string  `json:"institution"`
	Balance     float64 `json:"balance"`
}

func RegisterAccountRoutes(r *mux.Router) {
	sub := r.PathPrefix("/accounts").Subrouter()
	sub.Use(middleware.AuthMiddleware(models.DB))
	sub.HandleFunc("", getAccounts).Methods("GET")
	sub.HandleFunc("", createAccount).Methods("POST")
	sub.HandleFunc("/{id}", updateAccount).Methods("PUT")
	sub.HandleFunc("/{id}", deleteAccount).Methods("DELETE")
}

func getAccounts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userCtx := middleware.GetUserContext(ctx)
	if userCtx == nil {
		writeJSONError(w, "Authentication required", http.StatusUnauthorized)
		return
	}
	
	var accounts []models.Account
	if err := models.DB.WithContext(ctx).Where("household_id IN ?", userCtx.HouseholdIDs).Find(&accounts).Error; err != nil {
		writeJSONError(w, "Failed to retrieve accounts", http.StatusInternalServerError)
		return
	}
	
	writeJSON(w, accounts, http.StatusOK)
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userCtx := middleware.GetUserContext(ctx)
	if userCtx == nil {
		writeJSONError(w, "Authentication required", http.StatusUnauthorized)
		return
	}
	
	var req AccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate required fields
	if req.Name == "" || req.Type == "" {
		writeJSONError(w, "Missing required fields: name and type", http.StatusBadRequest)
		return
	}
	
	// Validate account type
	validTypes := map[string]bool{
		"checking":  true,
		"savings":   true,
		"credit":    true,
		"investment": true,
		"loan":      true,
	}
	if !validTypes[req.Type] {
		writeJSONError(w, "Invalid account type", http.StatusBadRequest)
		return
	}
	
	if len(userCtx.HouseholdIDs) == 0 {
		writeJSONError(w, "No household found for user", http.StatusBadRequest)
		return
	}
	
	account := models.Account{
		Name:        req.Name,
		Type:        req.Type,
		Institution: req.Institution,
		Balance:     req.Balance,
		HouseholdID: userCtx.HouseholdIDs[0],
	}
	
	if err := models.DB.WithContext(ctx).Create(&account).Error; err != nil {
		writeJSONError(w, "Failed to create account", http.StatusInternalServerError)
		return
	}
	
	writeJSONSuccess(w, "Account created successfully", account)
}

func updateAccount(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userCtx := middleware.GetUserContext(ctx)
	if userCtx == nil {
		writeJSONError(w, "Authentication required", http.StatusUnauthorized)
		return
	}
	
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		writeJSONError(w, "Invalid account ID", http.StatusBadRequest)
		return
	}
	
	if !middleware.CheckAccountOwnership(ctx, userCtx.ID, uint(id)) {
		writeJSONError(w, "Access denied", http.StatusForbidden)
		return
	}

	var req AccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	var acc models.Account
	if err := models.DB.WithContext(ctx).First(&acc, id).Error; err != nil {
		writeJSONError(w, "Account not found", http.StatusNotFound)
		return
	}
	
	// Update only allowed fields
	if req.Name != "" {
		acc.Name = req.Name
	}
	if req.Type != "" {
		validTypes := map[string]bool{
			"checking":  true,
			"savings":   true,
			"credit":    true,
			"investment": true,
			"loan":      true,
		}
		if !validTypes[req.Type] {
			writeJSONError(w, "Invalid account type", http.StatusBadRequest)
			return
		}
		acc.Type = req.Type
	}
	if req.Institution != "" {
		acc.Institution = req.Institution
	}
	acc.Balance = req.Balance // Allow zero balance
	
	if err := models.DB.WithContext(ctx).Save(&acc).Error; err != nil {
		writeJSONError(w, "Failed to update account", http.StatusInternalServerError)
		return
	}
	
	writeJSONSuccess(w, "Account updated successfully", acc)
}

func deleteAccount(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	userCtx := middleware.GetUserContext(ctx)
	if userCtx == nil {
		writeJSONError(w, "Authentication required", http.StatusUnauthorized)
		return
	}
	
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		writeJSONError(w, "Invalid account ID", http.StatusBadRequest)
		return
	}
	
	if !middleware.CheckAccountOwnership(ctx, userCtx.ID, uint(id)) {
		writeJSONError(w, "Access denied", http.StatusForbidden)
		return
	}
	
	var acc models.Account
	if err := models.DB.WithContext(ctx).First(&acc, id).Error; err != nil {
		writeJSONError(w, "Account not found", http.StatusNotFound)
		return
	}
	
	if err := models.DB.WithContext(ctx).Delete(&acc).Error; err != nil {
		writeJSONError(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}
	
	writeJSONSuccess(w, "Account deleted successfully", nil)
}