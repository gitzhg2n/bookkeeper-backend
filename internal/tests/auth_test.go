package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"bookkeeper-backend/config"
	"bookkeeper-backend/internal/db"
	"bookkeeper-backend/routes"

	"gorm.io/gorm"
)

type testEnv struct {
	DB     *gorm.DB
	Server http.Handler
}

func slogDiscard() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func setupTest(t *testing.T) *testEnv {
	t.Helper()
	_ = os.Setenv("JWT_SECRET", "0123456789abcdefghijklmnopqrstuvwxyz012345")
	cfg := config.Load()
	cfg.DatabaseURL = ":memory:"

	_, gdb, err := db.Initialize(cfg)
	if err != nil {
		t.Fatalf("init db: %v", err)
	}

	logger := slogDiscard()
	srv := routes.BuildRouter(cfg, gdb, logger)
	return &testEnv{DB: gdb, Server: srv}
}

func TestRegisterLogin(t *testing.T) {
	env := setupTest(t)

	// Register
	regBody := `{"email":"a@example.com","password":"VerySecurePass1!"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/v1/auth/register", bytes.NewBufferString(regBody))
	env.Server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}
	var regResp struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &regResp)
	if regResp.Data.AccessToken == "" {
		t.Fatalf("missing access token")
	}

	// Login
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("POST", "/v1/auth/login", bytes.NewBufferString(regBody))
	env.Server.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("login expected 200 got %d body=%s", w2.Code, w2.Body.String())
	}
}

func TestInvalidLogin(t *testing.T) {
	env := setupTest(t)

	// Attempt login before register
	body := `{"email":"nouser@example.com","password":"pw"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/v1/auth/login", bytes.NewBufferString(body))
	env.Server.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w.Code)
	}
}