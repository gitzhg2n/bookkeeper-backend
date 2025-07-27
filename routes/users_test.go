package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"bookkeeper-backend/models"
)

func TestCreateUser_Admin(t *testing.T) {
	user := models.User{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		RecoverySeedHash: "seedhash",
		Role:         "user",
	}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/users/", bytes.NewReader(body))
	ctx := req.Context()
	// Mock admin role in context
	ctx = contextWithRole(ctx, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	createUser(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Errorf("expected 200 or 201 for admin, got %d", w.Code)
	}
}

func TestCreateUser_NonAdmin(t *testing.T) {
	user := models.User{
		Email:        "test2@example.com",
		PasswordHash: "hashedpassword",
		RecoverySeedHash: "seedhash",
		Role:         "user",
	}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/users/", bytes.NewReader(body))
	ctx := req.Context()
	// Mock non-admin role in context
	ctx = contextWithRole(ctx, "user")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	createUser(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 for non-admin, got %d", w.Code)
	}
}

