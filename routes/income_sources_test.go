package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"bookkeeper-backend/models"
)

func TestCreateIncomeSource(t *testing.T) {
	src := models.IncomeSource{
		Name:        "Acme Corp Salary",
		Type:        "W-2",
		Amount:      5000,
		HouseholdID: 1,
		Frequency:   "monthly",
		Notes:       "Regular salary",
	}
	body, _ := json.Marshal(src)
	req := httptest.NewRequest("POST", "/incomeSources/", bytes.NewReader(body))
	req = req.WithContext(mockUserContext(req.Context()))
	w := httptest.NewRecorder()
	createIncomeSource(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Errorf("expected 200 or 201, got %d", w.Code)
	}
}

func mockUserContext(ctx interface{}) interface{} {
	return ctx
}