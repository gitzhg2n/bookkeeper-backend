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

	mux.HandleFunc("/v1/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSONSuccess(r, w, "ok", map[string]string{"status": "up"})
	})

	// Auth + rate limiting (10 req per 60s per IP)
	rateLimiter := middleware.NewRateLimiter()
	authRateLimit := rateLimiter.Limit(60000, 10)

	authHandler := NewAuthHandler(cfg, gdb, logger)
	mux.Handle("/v1/auth/register", authRateLimit(http.HandlerFunc(authHandler.Register)))
	mux.Handle("/v1/auth/login", authRateLimit(http.HandlerFunc(authHandler.Login)))
	mux.Handle("/v1/auth/refresh", authRateLimit(http.HandlerFunc(authHandler.Refresh)))
	mux.Handle("/v1/auth/logout", authRateLimit(http.HandlerFunc(authHandler.Logout)))

	userHandler := NewUserHandler(gdb)
	households := NewHouseholdHandler(gdb)
	accounts := NewAccountHandler(gdb)
	transactions := NewTransactionHandler(gdb)
	categories := NewCategoryHandler(gdb)
	budgets := NewBudgetHandler(gdb)

	protected := middleware.AuthMiddleware(cfg)

	mux.Handle("/v1/users/me", protected(http.HandlerFunc(userHandler.Me)))

	mux.Handle("/v1/households", protected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			households.Create(w, r)
		case http.MethodGet:
			households.List(w, r)
		default:
			writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/v1/households/", protected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/v1/households/")
		parts := strings.Split(path, "/")
		if len(parts) < 2 {
			writeJSONError(r, w, "not found", http.StatusNotFound)
			return
		}
		householdID := parts[0]
		switch parts[1] {
		case "accounts":
			switch r.Method {
			case http.MethodPost:
				accounts.Create(w, r, householdID)
			case http.MethodGet:
				accounts.List(w, r, householdID)
			default:
				writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		case "categories":
			switch r.Method {
			case http.MethodPost:
				categories.Create(w, r, householdID)
			case http.MethodGet:
				categories.List(w, r, householdID)
			default:
				writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		case "budgets":
			if len(parts) == 2 {
				switch r.Method {
				case http.MethodPost:
					budgets.Create(w, r, householdID)
				case http.MethodGet:
					budgets.List(w, r, householdID)
				case http.MethodPut:
					budgets.Upsert(w, r, householdID)
				default:
					writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
				}
				return
			}
			if len(parts) == 3 {
				if r.Method == http.MethodDelete {
					budgets.Delete(w, r, householdID, parts[2])
					return
				}
				writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			writeJSONError(r, w, "not found", http.StatusNotFound)
			return
		case "budget_summary":
			if r.Method == http.MethodGet {
				budgets.Summary(w, r, householdID)
				return
			}
			writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
			return
		default:
			writeJSONError(r, w, "not found", http.StatusNotFound)
			return
		}
	})))

	mux.Handle("/v1/accounts/", protected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/v1/accounts/")
		parts := strings.Split(path, "/")
		if len(parts) >= 2 && parts[1] == "transactions" {
			accountID := parts[0]
			switch r.Method {
			case http.MethodPost:
				transactions.Create(w, r, accountID)
			case http.MethodGet:
				transactions.List(w, r, accountID)
			default:
				writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}
		writeJSONError(r, w, "not found", http.StatusNotFound)
	})))

	root := middleware.RequestID()(middleware.Logging(logger)(mux))
	return root
}