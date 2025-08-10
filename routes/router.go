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

	// Rate limiter for auth endpoints (10 requests / 60s per IP)
	authRateLimit := middleware.NewRateLimiter().Limit(60000, 10)

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

	// Households root
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

	// Nested resources under /v1/households/{id}/...
	mux.Handle("/v1/households/", protected(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/v1/households/")
		parts := strings.Split(path, "/")

		if len(parts) >= 2 {
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
				// Handle both /v1/households/{id}/budgets and /v1/households/{id}/budgets/{budgetID}
				if len(parts) == 2 {
					// /v1/households/{id}/budgets (GET list, POST create, PUT upsert)
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
				} else if len(parts) == 3 {
					// /v1/households/{id}/budgets/{budgetID}
					budgetID := parts[2]
					switch r.Method {
					case http.MethodDelete:
						budgets.Delete(w, r, householdID, budgetID)
					default:
						writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
					}
				} else {
					writeJSONError(r, w, "not found", http.StatusNotFound)
				}
				return
			case "budget_summary":
				// GET /v1/households/{id}/budget_summary?month=YYYY-MM
				if r.Method == http.MethodGet {
					budgets.Summary(w, r, householdID)
					return
				}
				writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
		}
		writeJSONError(r, w, "not found", http.StatusNotFound)
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
			writeJSONError(r, w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSONError(r, w, "not found", http.StatusNotFound)
	})))

	root := middleware.RequestID()(middleware.Logging(logger)(mux))
	return root
}