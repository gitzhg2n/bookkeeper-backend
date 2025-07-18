package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"bookkeeper-backend-go/models"
	"bookkeeper-backend-go/middleware"
)

func RegisterIncomeSourceRoutes(r *mux.Router) {
	sub := r.PathPrefix("/incomeSources").Subrouter()
	sub.HandleFunc("/", createIncomeSource).Methods("POST")
	sub.HandleFunc("/", listIncomeSources).Methods("GET")
	sub.HandleFunc("/{id}", getIncomeSource).Methods("GET")
	sub.HandleFunc("/{id}", updateIncomeSource).Methods("PUT")
	sub.HandleFunc("/{id}", deleteIncomeSource).Methods("DELETE")
}

func createIncomeSource(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	var source models.IncomeSource
	if err := json.NewDecoder(r.Body).Decode(&source); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	if !middleware.CheckHouseholdOwnership(r.Context(), userID, source.HouseholdID) {
		http.Error(w, "Forbidden: Not your household", http.StatusForbidden)
		return
	}
	if err := models.DB.Create(&source).Error; err != nil {
		http.Error(w, "Error creating income source", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(source)
}

func listIncomeSources(w http.ResponseWriter, r *http.Request) {
	householdIDs := r.Context().Value("householdIDs").([]uint)
	var sources []models.IncomeSource
	if err := models.DB.Where("household_id IN (?)", householdIDs).Order("created_at DESC").Find(&sources).Error; err != nil {
		http.Error(w, "Error fetching income sources", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(sources)
}

func getIncomeSource(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var source models.IncomeSource
	if err := models.DB.First(&source, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if !middleware.CheckHouseholdOwnership(r.Context(), userID, source.HouseholdID) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	json.NewEncoder(w).Encode(source)
}

func updateIncomeSource(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var source models.IncomeSource
	if err := models.DB.First(&source, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if !middleware.CheckHouseholdOwnership(r.Context(), userID, source.HouseholdID) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	var payload models.IncomeSource
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	source.Name = payload.Name
	source.Type = payload.Type
	source.Amount = payload.Amount
	source.Frequency = payload.Frequency
	source.Notes = payload.Notes
	if err := models.DB.Save(&source).Error; err != nil {
		http.Error(w, "Error updating income source", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(source)
}

func deleteIncomeSource(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var source models.IncomeSource
	if err := models.DB.First(&source, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if !middleware.CheckHouseholdOwnership(r.Context(), userID, source.HouseholdID) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if err := models.DB.Delete(&source).Error; err != nil {
		http.Error(w, "Error deleting income source", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Income source deleted"})
}