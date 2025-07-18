package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"context"
)

func TestAdminOnlyMiddleware(t *testing.T) {
	handler := AdminOnly(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Test with admin role
	req1, _ := http.NewRequest("GET", "/", nil)
	ctx1 := context.WithValue(req1.Context(), "role", "admin")
	req1 = req1.WithContext(ctx1)
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req1)
	if rr1.Code != http.StatusOK {
		t.Errorf("Admin should be allowed, got status %d", rr1.Code)
	}

	// Test with user role
	req2, _ := http.NewRequest("GET", "/", nil)
	ctx2 := context.WithValue(req2.Context(), "role", "user")
	req2 = req2.WithContext(ctx2)
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusForbidden {
		t.Errorf("Non-admin should be forbidden, got status %d", rr2.Code)
	}
}