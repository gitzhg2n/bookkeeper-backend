package routes

import (
	"log/slog"
	"net/http"
	"strings"

	"bookkeeper-backend/config"
	"bookkeeper-backend/middleware"

	"gorm.io/gorm"
)

func BuildRouter(cfg *config.Config, gdb *gorm.DB, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("/v1/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSONSuccess(w, "ok", map[string]string{"status": "up"})
	})

	authHandler := NewAuthHandler(cfg, gdb, logger)
	mux.HandleFunc("/v1/auth/register", authHandler.Register)
	mux.HandleFunc("/v1/auth/login", authHandler.Login)
	mux.HandleFunc("/v1/auth/refresh", authHandler.Refresh)
	mux.HandleFunc("/v1/auth/logout", authHandler.Logout)

	userHandler := NewUserHandler(gdb)
	households := NewHouseholdHandler(gdb)
	accounts := NewAccountHandler(gdb)
	transactions := NewTransactionHandler(gdb)

	protected := middleware.AuthMiddleware(cfg)

	// Users
	mux.Handle("/v1/users/me", protected(http.HandlerFunc(userHandler.Me)))

	// Households
	mux.Handle("/v1/households", protected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			households.Create(w, r)
		case http.MethodGet:
			households.List(w, r)
		default:
			writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Pattern: /v1/households/{id}/accounts
	mux.Handle("/v1/households/", protected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/v1/households/")
		parts := strings.Split(path, "/")
		if len(parts) >= 2 && parts[1] == "accounts" {
			householdID := parts[0]
			switch r.Method {
			case http.MethodPost:
				accounts.Create(w, r, householdID)
				return
			case http.MethodGet:
				accounts.List(w, r, householdID)
				return
			}
			writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSONError(w, "not found", http.StatusNotFound)
	})))

	// Accounts transactions: /v1/accounts/{id}/transactions
	mux.Handle("/v1/accounts/", protected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/v1/accounts/")
		parts := strings.Split(path, "/")
		if len(parts) >= 2 && parts[1] == "transactions" {
			accountID := parts[0]
			switch r.Method {
			case http.MethodPost:
				transactions.Create(w, r, accountID)
				return
			case http.MethodGet:
				transactions.List(w, r, accountID)
				return
			}
			writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSONError(w, "not found", http.StatusNotFound)
	})))

	// Wrap with middleware: request id + logging
	root := middleware.RequestID()(middleware.Logging(logger)(mux))
	return root
}