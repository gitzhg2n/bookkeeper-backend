package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHouseholdIsolation(t *testing.T) {
	env := setupTest(t)

	// User A
	regA := `{"email":"ha@example.com","password":"StrongPassw0rd!"}`
	w := httptest.NewRecorder()
	env.Server.ServeHTTP(w, httptest.NewRequest("POST", "/v1/auth/register", bytes.NewBufferString(regA)))
	if w.Code != 200 {
		t.Fatalf("reg A failed: %d %s", w.Code, w.Body.String())
	}
	tokenA := extractToken(t, w.Body.Bytes())

	// User B
	regB := `{"email":"hb@example.com","password":"StrongPassw0rd!"}`
	w = httptest.NewRecorder()
	env.Server.ServeHTTP(w, httptest.NewRequest("POST", "/v1/auth/register", bytes.NewBufferString(regB)))
	if w.Code != 200 {
		t.Fatalf("reg B failed")
	}
	tokenB := extractToken(t, w.Body.Bytes())

	// A creates household
	w = httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/v1/households", bytes.NewBufferString(`{"name":"AHome"}`))
	req.Header.Set("Authorization", "Bearer "+tokenA)
	env.Server.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("household create A failed: %d %s", w.Code, w.Body.String())
	}
	hID := extractID(t, w.Body.Bytes())

	// B tries to list accounts of A's household
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", fmt.Sprintf("/v1/households/%d/accounts", hID), nil)
	req.Header.Set("Authorization", "Bearer "+tokenB)
	env.Server.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected forbidden for cross-household, got %d", w.Code)
	}
}