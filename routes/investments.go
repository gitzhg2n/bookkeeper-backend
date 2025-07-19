package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"bookkeeper-backend-go/models"
	"github.com/gorilla/mux"
	"bookkeeper-backend-go/middleware"
)

func RegisterInvestmentRoutes(r *mux.Router) {
	sub := r.PathPrefix("/investments").Subrouter()
	sub.HandleFunc("", getInvestments).Methods("GET")
	sub.HandleFunc("", createInvestment).Methods("POST")
	sub.HandleFunc("/{id}", updateInvestment).Methods("PUT")
	sub.HandleFunc("/{id}", deleteInvestment).Methods("DELETE")
}

type InvestmentRequest struct {
	Name        string  `json:"name"`
	Value       float64 `json:"value"`
	Type        string  `json:"type"`
	Institution string  `json:"institution"`
}

func getInvestments(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserContext(r.Context())
	var investments []models.Investment
	models.DB.Where("user_id = ?", user.ID).Find(&investments)
	json.NewEncoder(w).Encode(investments)
}

func createInvestment(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserContext(r.Context())
	var req InvestmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Type == "" || req.Value < 0 || req.Institution == "" {
		http.Error(w, "Missing/invalid fields", http.StatusBadRequest)
		return
	}
	investment := models.Investment{
		UserID:      user.ID,
		Name:        req.Name,
		Value:       req.Value,
		Type:        req.Type,
		Institution: req.Institution,
	}
	if err := models.DB.Create(&investment).Error; err != nil {
		http.Error(w, "Failed to create investment", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(investment)
}

func updateInvestment(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserContext(r.Context())
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var investment models.Investment
	if err := models.DB.First(&investment, id).Error; err != nil {
		http.Error(w, "Investment not found", http.StatusNotFound)
		return
	}
	if investment.UserID != user.ID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	var req InvestmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Type == "" || req.Value < 0 || req.Institution == "" {
		http.Error(w, "Missing/invalid fields", http.StatusBadRequest)
		return
	}
	investment.Name = req.Name
	investment.Value = req.Value
	investment.Type = req.Type
	investment.Institution = req.Institution
	if err := models.DB.Save(&investment).Error; err != nil {
		http.Error(w, "Failed to update investment", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(investment)
}

func deleteInvestment(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserContext(r.Context())
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var investment models.Investment
	if err := models.DB.First(&investment, id).Error; err != nil {
		http.Error(w, "Investment not found", http.StatusNotFound)
		return
	}
	if investment.UserID != user.ID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if err := models.DB.Delete(&investment).Error; err != nil {
		http.Error(w, "Failed to delete investment", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Investment deleted"})
}