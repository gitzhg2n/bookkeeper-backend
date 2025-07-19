package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"bookkeeper-backend/models"
)

func TestCreateInvestment(t *testing.T) {
	inv := models.Investment{
		Name:        "Vanguard S&P 500",
		AccountID:   1,
		Type:        "mutual fund",
		Institution: "Vanguard",
	}
	body, _ := json.Marshal(inv)
	req := httptest.NewRequest("POST", "/investments/", bytes.NewReader(body))
	req = req.WithContext(mockUserContext(req.Context()))
	w := httptest.NewRecorder()
	createInvestment(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Errorf("expected 200 or 201, got %d", w.Code)
	}
}

func mockUserContext(ctx interface{}) interface{} {
	return ctx
}