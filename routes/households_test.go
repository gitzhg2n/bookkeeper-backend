package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"bookkeeper-backend/models"
)

func TestCreateHousehold(t *testing.T) {
	household := models.Household{
		Name: "Smith Family",
		OwnerID: 42,
	}
	body, _ := json.Marshal(household)
	req := httptest.NewRequest("POST", "/households/", bytes.NewReader(body))
	ctx := req.Context()
	ctx = contextWithUserID(ctx, 42)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	createHousehold(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Errorf("expected 200 or 201, got %d", w.Code)
	}
}

// Helper to mock userID in context
func contextWithUserID(ctx interface{}, userID uint) interface{} {
	// Extend as needed for your context infrastructure
	return ctx
}