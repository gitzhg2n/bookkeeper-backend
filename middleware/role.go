package middleware

import (
	"net/http"
)

// AdminOnly middleware ensures only admin users can access the endpoint
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userCtx := GetUserContext(r.Context())
		if userCtx == nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}
		
		if userCtx.Role != "admin" {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// RequireRole creates middleware that requires a specific role
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userCtx := GetUserContext(r.Context())
			if userCtx == nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
			
			if userCtx.Role != role {
				http.Error(w, "Insufficient privileges", http.StatusForbidden)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
}