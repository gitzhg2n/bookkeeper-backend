package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"bookkeeper-backend/models"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type registerPayload struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	RecoverySeed  string `json:"recoverySeed"`
}

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type recoverPayload struct {
	Email        string `json:"email"`
	RecoverySeed string `json:"recoverySeed"`
	NewPassword  string `json:"newPassword"`
}

func RegisterAuthRoutes(r *mux.Router) {
	sub := r.PathPrefix("/auth").Subrouter()
	sub.HandleFunc("/register", registerUser).Methods("POST")
	sub.HandleFunc("/login", loginUser).Methods("POST")
	sub.HandleFunc("/recover", recoverUser).Methods("POST")
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	var payload registerPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}
	if payload.Email == "" || payload.Password == "" || payload.RecoverySeed == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}
	var count int64
	models.DB.Model(&models.User{}).Where("email = ?", payload.Email).Count(&count)
	if count > 0 {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), 12)
	recoverySeedHash := models.HashRecoverySeed(payload.RecoverySeed)
	user := models.User{
		Email:            payload.Email,
		PasswordHash:     string(passwordHash),
		RecoverySeedHash: recoverySeedHash,
	}
	if err := models.DB.Create(&user).Error; err != nil {
		http.Error(w, "Server error during registration", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered"})
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var payload loginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}
	var user models.User
	if err := models.DB.Where("email = ?", payload.Email).First(&user).Error; err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	tokenString, _ := token.SignedString(jwtSecret)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func recoverUser(w http.ResponseWriter, r *http.Request) {
	var payload recoverPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}
	var user models.User
	if err := models.DB.Where("email = ?", payload.Email).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if !models.CheckRecoverySeed(payload.RecoverySeed, user.RecoverySeedHash) {
		http.Error(w, "Invalid recovery seed", http.StatusUnauthorized)
		return
	}
	newHash, _ := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), 12)
	user.PasswordHash = string(newHash)
	if err := models.DB.Save(&user).Error; err != nil {
		http.Error(w, "Server error during recovery", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Password reset successful"})
}