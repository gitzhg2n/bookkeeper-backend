package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"bookkeeper-backend-go/models"
	"bookkeeper-backend-go/middleware"
)

func RegisterUserRoutes(r *mux.Router) {
	sub := r.PathPrefix("/users").Subrouter()
	sub.HandleFunc("/", createUser).Methods("POST")
	sub.HandleFunc("/", listUsers).Methods("GET")
	sub.HandleFunc("/{id}", getUser).Methods("GET")
	sub.HandleFunc("/{id}", updateUser).Methods("PUT")
	sub.HandleFunc("/{id}", deleteUser).Methods("DELETE")
}

func createUser(w http.ResponseWriter, r *http.Request) {
	// Only admin can create users
	role := r.Context().Value("role").(string)
	if role != "admin" {
		http.Error(w, "Admin access required", http.StatusForbidden)
		return
	}
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	if err := models.DB.Create(&user).Error; err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func listUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := models.DB.Select("id", "email").Find(&users).Error; err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var user models.User
	if err := models.DB.Select("id", "email").First(&user, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	// Only admin can update users
	role := r.Context().Value("role").(string)
	if role != "admin" {
		http.Error(w, "Admin access required", http.StatusForbidden)
		return
	}
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var user models.User
	if err := models.DB.First(&user, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	var payload models.User
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	user.Email = payload.Email
	if err := models.DB.Save(&user).Error; err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	// Only admin can delete users
	role := r.Context().Value("role").(string)
	if role != "admin" {
		http.Error(w, "Admin access required", http.StatusForbidden)
		return
	}
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var user models.User
	if err := models.DB.First(&user, id).Error; err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if err := models.DB.Delete(&user).Error; err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted"})
}