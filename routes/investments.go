package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"bookkeeper-backend-go/models"
	"bookkeeper-backend-go/middleware"
)

func RegisterInvestmentRoutes(r *mux.Router) {
	sub := r.PathPrefix("/investments").Subrouter()
	sub.HandleFunc("/", createInvestment).Methods("POST")
	sub.HandleFunc("/", listInvestments).Methods("GET")
	sub.HandleFunc("/{id}", getInvestment).Methods("GET")
	sub.HandleFunc("/{id}", updateInvestment).Methods("PUT")
	sub.HandleFunc("/{id}", deleteInvestment).Methods("DELETE")
}

func createInvestment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	var inv models.Investment
	if err := json.NewDecoder(r.Body).Decode(&inv); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	if !middleware.CheckAccountOwnership(r.Context(), userID, inv.AccountID) {
		http.Error(w, "Forbidden: Not your account", http.StatusForbidden)
		return
	}
	if err := models.DB.Create(&inv).Error; err != nil {
		http.Error(w, "Error creating investment", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(inv)
}

func listInvestments(w http.ResponseWriter, r *http.Request) {
	accountIDs := r.Context().Value("accountIDs").([]uint)
	var investments []models.Investment
	if err := models.DB.Where("account_id IN (?)", accountIDs).Find(&investments).Error; err != nil {
		http.Error(w, "Error fetching investments", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(investments)
}

func getInvestment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var inv models.Investment
	if err := models.DB.First(&inv, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if !middleware.CheckAccountOwnership(r.Context(), userID, inv.AccountID) {
		http.Error(w, "Forbidden: Not your investment", http.StatusForbidden)
		return
	}
	json.NewEncoder(w).Encode(inv)
}

func updateInvestment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var inv models.Investment
	if err := models.DB.First(&inv, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if !middleware.CheckAccountOwnership(r.Context(), userID, inv.AccountID) {
		http.Error(w, "Forbidden: Not your investment", http.StatusForbidden)
		return
	}
	var payload models.Investment
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	inv.Name = payload.Name
	inv.Type = payload.Type
	inv.Institution = payload.Institution
	if err := models.DB.Save(&inv).Error; err != nil {
		http.Error(w, "Error updating investment", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(inv)
}

func deleteInvestment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var inv models.Investment
	if err := models.DB.First(&inv, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if !middleware.CheckAccountOwnership(r.Context(), userID, inv.AccountID) {
		http.Error(w, "Forbidden: Not your investment", http.StatusForbidden)
		return
	}
	if err := models.DB.Delete(&inv).Error; err != nil {
		http.Error(w, "Error deleting investment", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Investment deleted"})
}