package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"bookkeeper-backend/models"
	"bookkeeper-backend/middleware"
)

func RegisterHouseholdRoutes(r *mux.Router) {
	sub := r.PathPrefix("/households").Subrouter()
	sub.HandleFunc("/", createHousehold).Methods("POST")
	sub.HandleFunc("/", listHouseholds).Methods("GET")
	sub.HandleFunc("/{id}", getHousehold).Methods("GET")
	sub.HandleFunc("/{id}", updateHousehold).Methods("PUT")
	sub.HandleFunc("/{id}", deleteHousehold).Methods("DELETE")
}

func createHousehold(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	var household models.Household
	if err := json.NewDecoder(r.Body).Decode(&household); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	household.OwnerID = userID
	if err := models.DB.Create(&household).Error; err != nil {
		http.Error(w, "Error creating household", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(household)
}

func listHouseholds(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	var households []models.Household
	if err := models.DB.Where("owner_id = ?", userID).Find(&households).Error; err != nil {
		http.Error(w, "Error fetching households", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(households)
}

func getHousehold(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var household models.Household
	if err := models.DB.First(&household, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if !middleware.CheckHouseholdOwnership(r.Context(), userID, household.ID) {
		http.Error(w, "Forbidden: Not your household", http.StatusForbidden)
		return
	}
	json.NewEncoder(w).Encode(household)
}

func updateHousehold(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var household models.Household
	if err := models.DB.First(&household, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if !middleware.CheckHouseholdOwnership(r.Context(), userID, household.ID) {
		http.Error(w, "Forbidden: Not your household", http.StatusForbidden)
		return
	}
	var payload models.Household
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	household.Name = payload.Name
	if err := models.DB.Save(&household).Error; err != nil {
		http.Error(w, "Error updating household", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(household)
}

func deleteHousehold(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var household models.Household
	if err := models.DB.First(&household, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if !middleware.CheckHouseholdOwnership(r.Context(), userID, household.ID) {
		http.Error(w, "Forbidden: Not your household", http.StatusForbidden)
		return
	}
	if err := models.DB.Delete(&household).Error; err != nil {
		http.Error(w, "Error deleting household", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Household deleted"})
}