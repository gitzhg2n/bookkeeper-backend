package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"bookkeeper-backend-go/models"
	"bookkeeper-backend-go/middleware"
)

func RegisterAccountRoutes(r *mux.Router) {
	sub := r.PathPrefix("/accounts").Subrouter()
	sub.HandleFunc("/", createAccount).Methods("POST")
	sub.HandleFunc("/", listAccounts).Methods("GET")
	sub.HandleFunc("/{id}", getAccount).Methods("GET")
	sub.HandleFunc("/{id}", updateAccount).Methods("PUT")
	sub.HandleFunc("/{id}", deleteAccount).Methods("DELETE")
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	var acc models.Account
	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	if !middleware.CheckHouseholdOwnership(r.Context(), userID, acc.HouseholdID) {
		http.Error(w, "Forbidden: Not your household", http.StatusForbidden)
		return
	}
	if err := models.DB.Create(&acc).Error; err != nil {
		http.Error(w, "Error creating account", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(acc)
}

func listAccounts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	householdIDs := r.Context().Value("householdIDs").([]uint)
	var accounts []models.Account
	if err := models.DB.Where("household_id IN (?)", householdIDs).Find(&accounts).Error; err != nil {
		http.Error(w, "Error fetching accounts", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(accounts)
}

func getAccount(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var acc models.Account
	if err := models.DB.First(&acc, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if !middleware.CheckHouseholdOwnership(r.Context(), userID, acc.HouseholdID) {
		http.Error(w, "Forbidden: Not your account", http.StatusForbidden)
		return
	}
	json.NewEncoder(w).Encode(acc)
}

func updateAccount(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var acc models.Account
	if err := models.DB.First(&acc, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if !middleware.CheckHouseholdOwnership(r.Context(), userID, acc.HouseholdID) {
		http.Error(w, "Forbidden: Not your account", http.StatusForbidden)
		return
	}
	var payload models.Account
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	// Only update allowed fields
	acc.Name = payload.Name
	acc.Type = payload.Type
	acc.Institution = payload.Institution
	if err := models.DB.Save(&acc).Error; err != nil {
		http.Error(w, "Error updating account", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(acc)
}

func deleteAccount(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var acc models.Account
	if err := models.DB.First(&acc, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if !middleware.CheckHouseholdOwnership(r.Context(), userID, acc.HouseholdID) {
		http.Error(w, "Forbidden: Not your account", http.StatusForbidden)
		return
	}
	if err := models.DB.Delete(&acc).Error; err != nil {
		http.Error(w, "Error deleting account", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Account deleted"})
}