package middleware

import (
	"context"
	"net/http"
	"strings"

	"bookkeeper-backend-go/models"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type UserContext struct {
	ID           uint
	Email        string
	Role         string
	HouseholdIDs []uint
	AccountIDs   []uint
}

// AuthMiddleware verifies JWT and enriches request context with user info.
func AuthMiddleware
