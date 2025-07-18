package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	rl := NewRateLimiter()
	handler := rl.Limit(1000, 2)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, _ := http.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1:1234"

	// First request: should pass
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req)
	if rr1.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", rr1.Code)
	}

	// Second request: should pass
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req)
	if rr2.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", rr2.Code)
	}

	// Third request: should fail (rate limited)
	rr3 := httptest.NewRecorder()
	handler.ServeHTTP(rr3, req)
	if rr3.Code != http.StatusTooManyRequests {
		t.Errorf("Expected 429, got %d", rr3.Code)
	}

	// Wait for window to expire, should pass again
	time.Sleep(1100 * time.Millisecond)
	rr4 := httptest.NewRecorder()
	handler.ServeHTTP(rr4, req)
	if rr4.Code != http.StatusOK {
		t.Errorf("Expected 200 after window reset, got %d", rr4.Code)
	}
}