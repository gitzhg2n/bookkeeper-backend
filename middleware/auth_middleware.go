package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"bookkeeper-backend/models"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
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

// UserContext holds info for privacy/ownership checks.
type UserContext struct {
	ID           uint
	Email        string
	Role         string
	HouseholdIDs []uint
	AccountIDs   []uint
}

// contextKey is a custom type for context values to avoid collisions.
type contextKey string

var userContextKey = contextKey("userContext")

// AuthMiddleware verifies JWT, loads user data, and attaches context.
func AuthMiddleware(db *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract JWT from Authorization header.
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Parse and verify JWT with signature.
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return jwtSecret, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Extract claims (assuming custom claims with user_id, email, role).
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
				return
			}
			userID := uint(userIDFloat)

			email, ok := claims["email"].(string)
			if !ok {
				http.Error(w, "Invalid email in token", http.StatusUnauthorized)
				return
			}

			role, ok := claims["role"].(string)
			if !ok {
				http.Error(w, "Invalid role in token", http.StatusUnauthorized)
				return
			}

			// Load user from DB to verify existence.
			var user models.User
			if err := db.First(&user, userID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					http.Error(w, "User not found", http.StatusUnauthorized)
				} else {
					http.Error(w, "Database error", http.StatusInternalServerError)
				}
				return
			}

			// Load household IDs (assuming User has many Households).
			var households []models.Household // Adjust to your model name.
			if err := db.Model(&user).Association("Households").Find(&households); err != nil {
				http.Error(w, "Failed to load households", http.StatusInternalServerError)
				return
			}
			householdIDs := make([]uint, len(households))
			for i, h := range households {
				householdIDs[i] = h.ID // Assuming ID field.
			}

			// Load account IDs (assuming User has many Accounts, or via Households—adjust as needed).
			var accounts []models.Account // Adjust to your model name.
			if err := db.Model(&user).Association("Accounts").Find(&accounts); err != nil {
				http.Error(w, "Failed to load accounts", http.StatusInternalServerError)
				return
			}
			accountIDs := make([]uint, len(accounts))
			for i, a := range accounts {
				accountIDs[i] = a.ID // Assuming ID field.
			}

			// Create UserContext and attach to request context.
			userCtx := UserContext{
				ID:           user.ID,
				Email:        email,
				Role:         role,
				HouseholdIDs: householdIDs,
				AccountIDs:   accountIDs,
			}
			ctx := context.WithValue(r.Context(), userContextKey, &userCtx)

			// Proceed to next handler with updated context.
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserContext retrieves UserContext from context
func GetUserContext(ctx context.Context) *UserContext {
	val := ctx.Value(userContextKey)
	if val == nil {
		return nil
	}
	return val.(*UserContext)
}
