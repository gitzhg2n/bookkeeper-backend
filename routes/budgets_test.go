package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"bookkeeper-backend-go/models"
)

func TestCreateBudget(t *testing.T) {
	budget := models.Budget{
		Name:        "Groceries",
		HouseholdID: 1,
		Period:      "monthly",
		Category:    "Food",
	}
	body, _ := json.Marshal(budget)
	req := httptest.NewRequest("POST", "/budgets/", bytes.NewReader(body))
	req = req.WithContext(mockUserContext(req.Context()))
	w := httptest.NewRecorder()
	createBudget(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Errorf("expected 200 or 201, got %d", w.Code)
	}
}