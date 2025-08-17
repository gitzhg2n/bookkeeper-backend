package routes

import (
	"log/slog"
	"net/http"
	"strings"

	"bookkeeper-backend/config"
	"bookkeeper-backend/internal/db"
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

	authHandler := NewAuthHandler(cfg, gdb, logger, &notificationStore)
	mux.Handle("/v1/auth/register", authRateLimit(http.HandlerFunc(authHandler.Register)))
	mux.Handle("/v1/auth/login", authRateLimit(http.HandlerFunc(authHandler.Login)))
	mux.Handle("/v1/auth/refresh", authRateLimit(http.HandlerFunc(authHandler.Refresh)))
	mux.Handle("/v1/auth/logout", authRateLimit(http.HandlerFunc(authHandler.Logout)))

	userHandler := NewUserHandler(gdb)
	households := NewHouseholdHandler(gdb)
	accounts := NewAccountHandler(gdb)
	transactions := NewTransactionHandler(gdb)
	categories := NewCategoryHandler(gdb)
	notificationStore := db.NotificationStore{DB: gdb.DB()}
	budgets := NewBudgetHandler(gdb, &notificationStore)
	// calculators are implemented as package-level handlers

	protected := middleware.AuthMiddleware(cfg)

	// admin entitlement management (admin-only endpoints)
	adminEnt := NewAdminEntitlementHandler(gdb)
	mux.Handle("/v1/admin/entitlements", protected(http.HandlerFunc(adminEnt.List)))
	mux.Handle("/v1/admin/entitlements/upsert", protected(http.HandlerFunc(adminEnt.Upsert)))

	// Calculators
	mux.Handle("/v1/calculators/mortgage", protected(http.HandlerFunc(MortgageCalculator)))
	mux.Handle("/v1/calculators/debt-payoff", protected(http.HandlerFunc(DebtPayoffCalculator)))
	mux.Handle("/v1/calculators/investment-growth", protected(http.HandlerFunc(InvestmentGrowthCalculator)))
	mux.Handle("/v1/calculators/rent-vs-buy", protected(http.HandlerFunc(RentVsBuyCalculator)))
	mux.Handle("/v1/calculators/tax-estimator", protected(http.HandlerFunc(TaxEstimatorCalculator)))

	mux.Handle("/v1/calculators/amortization", protected(http.HandlerFunc(AmortizationScheduleHandler)))
	mux.Handle("/v1/calculators/refinance-breakeven", protected(http.HandlerFunc(RefinanceBreakevenHandler)))
	mux.Handle("/v1/calculators/apr-to-apy", protected(http.HandlerFunc(APRToAPYHandler)))
	mux.Handle("/v1/calculators/apy-to-apr", protected(http.HandlerFunc(APYToAPRHandler)))
	mux.Handle("/v1/calculators/retirement-projection", protected(http.HandlerFunc(RetirementProjectionHandler)))
	mux.Handle("/v1/calculators/savings-goal", protected(http.HandlerFunc(SavingsGoalHandler)))
	mux.Handle("/v1/calculators/credit-payoff", protected(http.HandlerFunc(CreditPayoffHandler)))
	mux.Handle("/v1/calculators/take-home", protected(http.HandlerFunc(TakeHomeHandler)))
	mux.Handle("/v1/calculators/inflation-adjust", protected(http.HandlerFunc(InflationAdjustHandler)))
	mux.Handle("/v1/calculators/net-worth", protected(http.HandlerFunc(NetWorthHandler)))

	mux.Handle("/v1/calculators/loan-comparison", protected(http.HandlerFunc(LoanComparisonHandler)))
	mux.Handle("/v1/calculators/affordability", protected(http.HandlerFunc(AffordabilityHandler)))
	mux.Handle("/v1/calculators/credit-optimization", protected(http.HandlerFunc(CreditOptHandler)))
	mux.Handle("/v1/calculators/college-savings", protected(http.HandlerFunc(CollegeSavingsHandler)))
	mux.Handle("/v1/calculators/fee-drag", protected(http.HandlerFunc(FeeDragHandler)))
	mux.Handle("/v1/calculators/safe-withdrawal", protected(http.HandlerFunc(SafeWithdrawalHandler)))
	mux.Handle("/v1/calculators/cd-ladder", protected(http.HandlerFunc(CDLadderHandler)))
	mux.Handle("/v1/calculators/payroll", protected(http.HandlerFunc(PayrollHandler)))
	mux.Handle("/v1/calculators/convert-currency", protected(http.HandlerFunc(ConvertCurrencyHandler)))

	// Notifications
	notificationStore := db.NotificationStore{DB: gdb.DB()}
	notificationHandler := &NotificationHandler{Store: &notificationStore}
	mux.Handle("/v1/notifications", protected(http.HandlerFunc(notificationHandler.ListNotifications)))
	mux.Handle("/v1/notifications/read", protected(http.HandlerFunc(notificationHandler.MarkNotificationRead)))
	mux.Handle("/v1/notifications/read-all", protected(http.HandlerFunc(notificationHandler.MarkAllNotificationsRead)))

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

	userSettingsStore := db.UserSettingsStore{DB: gdb.DB()}
	userSettingsHandler := &UserSettingsHandler{Store: &userSettingsStore}
	mux.Handle("/v1/user/settings", protected(http.HandlerFunc(userSettingsHandler.Get)))
	mux.Handle("/v1/user/settings/update", protected(http.HandlerFunc(userSettingsHandler.Upsert)))

	root := middleware.RequestID()(middleware.Logging(logger)(mux))
	return root
}