package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"bookkeeper-backend/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	if len(secret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters long")
	}
	jwtSecret = []byte(secret)
}

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

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// isValidEmail validates email format
func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
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
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if payload.Email == "" || payload.Password == "" || payload.RecoverySeed == "" {
		writeJSONError(w, "Missing required fields: email, password, and recoverySeed", http.StatusBadRequest)
		return
	}
	
	// Basic email validation
	if !isValidEmail(payload.Email) {
		writeJSONError(w, "Invalid email format", http.StatusBadRequest)
		return
	}
	
	// Password strength validation
	if len(payload.Password) < 8 {
		writeJSONError(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}
	
	var count int64
	models.DB.Model(&models.User{}).Where("email = ?", payload.Email).Count(&count)
	if count > 0 {
		writeJSONError(w, "Email already registered", http.StatusConflict)
		return
	}
	
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 12)
	if err != nil {
		writeJSONError(w, "Server error during registration", http.StatusInternalServerError)
		return
	}
	recoverySeedHash, err := models.HashRecoverySeed(payload.RecoverySeed)
	if err != nil {
		writeJSONError(w, "Server error during registration", http.StatusInternalServerError)
		return
	}
	
	user := models.User{
		Email:            payload.Email,
		PasswordHash:     string(passwordHash),
		RecoverySeedHash: recoverySeedHash,
	}
	if err := models.DB.Create(&user).Error; err != nil {
		writeJSONError(w, "Server error during registration", http.StatusInternalServerError)
		return
	}
	
	writeJSONSuccess(w, "User registered successfully", map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	})
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var payload loginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if payload.Email == "" || payload.Password == "" {
		writeJSONError(w, "Missing required fields: email and password", http.StatusBadRequest)
		return
	}
	
	var user models.User
	if err := models.DB.Where("email = ?", payload.Email).First(&user).Error; err != nil {
		writeJSONError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)); err != nil {
		writeJSONError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat":     time.Now().Unix(),
	})
	
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		writeJSONError(w, "Server error during login", http.StatusInternalServerError)
		return
	}
	
	writeJSONSuccess(w, "Login successful", map[string]interface{}{
		"token": tokenString,
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func recoverUser(w http.ResponseWriter, r *http.Request) {
	var payload recoverPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if payload.Email == "" || payload.RecoverySeed == "" || payload.NewPassword == "" {
		writeJSONError(w, "Missing required fields: email, recoverySeed, and newPassword", http.StatusBadRequest)
		return
	}
	
	if len(payload.NewPassword) < 8 {
		writeJSONError(w, "New password must be at least 8 characters long", http.StatusBadRequest)
		return
	}
	
	var user models.User
	if err := models.DB.Where("email = ?", payload.Email).First(&user).Error; err != nil {
		writeJSONError(w, "User not found", http.StatusNotFound)
		return
	}
	
	if !models.CheckRecoverySeed(payload.RecoverySeed, user.RecoverySeedHash) {
		writeJSONError(w, "Invalid recovery seed", http.StatusUnauthorized)
		return
	}
	
	newHash, err := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), 12)
	if err != nil {
		writeJSONError(w, "Server error during recovery", http.StatusInternalServerError)
		return
	}
	
	user.PasswordHash = string(newHash)
	if err := models.DB.Save(&user).Error; err != nil {
		writeJSONError(w, "Server error during recovery", http.StatusInternalServerError)
		return
	}
	
	writeJSONSuccess(w, "Password reset successful", nil)
}