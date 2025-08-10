package routes

import (
	"net/http"

	"bookkeeper-backend/config"
	"bookkeeper-backend/middleware"

	"gorm.io/gorm"
)

func BuildRouter(cfg *config.Config, gdb *gorm.DB) http.Handler {
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("/v1/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSONSuccess(w, "ok", map[string]string{"status": "up"})
	})

	// Auth
	authHandler := NewAuthHandler(cfg, gdb)
	mux.HandleFunc("/v1/auth/register", authHandler.Register)
	mux.HandleFunc("/v1/auth/login", authHandler.Login)

	// Users (protected)
	userHandler := NewUserHandler(gdb)
	mux.Handle("/v1/users/me", middleware.AuthMiddleware(cfg)(http.HandlerFunc(userHandler.Me)))

	return mux
}