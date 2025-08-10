package middleware

import (
	"net/http"
	"strings"
	"time"

	"bookkeeper-backend/config"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"uid"`
	Email  string `json:"em"`
	Role   string `json:"r"`
	jwt.RegisteredClaims
}

func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authz := r.Header.Get("Authorization")
			if authz == "" || !strings.HasPrefix(authz, "Bearer ") {
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}
			tokenString := strings.TrimPrefix(authz, "Bearer ")

			token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				return cfg.JWTSecret, nil
			}, jwt.WithValidMethods([]string{"HS256"}))
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			claims, ok := token.Claims.(*Claims)
			if !ok || !token.Valid {
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}
			if claims.ExpiresAt == nil || time.Until(claims.ExpiresAt.Time) <= 0 {
				http.Error(w, "token expired", http.StatusUnauthorized)
				return
			}

			ctx := WithUser(r.Context(), &UserContext{
				ID:    claims.UserID,
				Email: claims.Email,
				Role:  claims.Role,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}