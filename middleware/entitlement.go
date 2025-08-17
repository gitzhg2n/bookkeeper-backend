package middleware

import (
    "encoding/json"
    "net/http"
    "strings"

    "bookkeeper-backend/config"
    "bookkeeper-backend/internal/db"
    "gorm.io/gorm"
)

// EntitlementMiddleware checks whether a user has access to a named feature.
// In soft-mode (hard=false) it sets header X-Feature-Allowed: true/false and allows the request (telemetry).
// In hard-mode it returns 402 Payment Required for denied access.
func EntitlementMiddleware(cfg *config.Config, gdb *gorm.DB, featureKey string, hard bool) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        fn := func(w http.ResponseWriter, r *http.Request) {
            // self-hosted gets full access
            if strings.ToLower(cfg.DeploymentMode) == "self-hosted" || strings.ToLower(cfg.DeploymentMode) == "selfhost" {
                w.Header().Set("X-Feature-Allowed", "true")
                next.ServeHTTP(w, r)
                return
            }

            user, ok := UserFrom(r.Context())
            if !ok {
                http.Error(w, "missing user context", http.StatusUnauthorized)
                return
            }

            allowed, err := db.UserHasEntitlement(gdb, user.ID, featureKey)
            if err != nil {
                w.Header().Set("X-Feature-Allowed", "false")
                if hard {
                    w.Header().Set("Content-Type", "application/json")
                    w.WriteHeader(http.StatusPaymentRequired)
                    json.NewEncoder(w).Encode(map[string]interface{}{"error": "feature_check_failed", "upgrade_url": "/pricing"})
                    return
                }
                next.ServeHTTP(w, r)
                return
            }

            if allowed {
                w.Header().Set("X-Feature-Allowed", "true")
                next.ServeHTTP(w, r)
                return
            }

            w.Header().Set("X-Feature-Allowed", "false")
            if hard {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusPaymentRequired)
                json.NewEncoder(w).Encode(map[string]interface{}{"error": "feature_not_available", "upgrade_url": "/pricing"})
                return
            }

            // soft-mode: allow but mark header for telemetry
            next.ServeHTTP(w, r)
        }
        return http.HandlerFunc(fn)
    }
}
