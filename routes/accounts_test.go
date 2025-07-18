package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"bookkeeper-backend-go/models"
)

func TestCreateAccount(t *testing.T) {
	acc := models.Account{
		Name:        "Checking",
		Type:        "checking",
		HouseholdID: 1,
		Institution: "Bank",
	}
	body, _ := json.Marshal(acc)
	req := httptest.NewRequest("POST", "/accounts/", bytes.NewReader(body))
	req = req.WithContext(mockUserContext(req.Context()))
	w := httptest.NewRecorder()
	createAccount(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Errorf("expected 200 or 201, got %d", w.Code)
	}
}

func mockUserContext(ctx interface{}) interface{} {
	// Extend this to mock userID, householdIDs, etc. as needed
	// For now, just return the context as-is for compilation
	return ctx
}