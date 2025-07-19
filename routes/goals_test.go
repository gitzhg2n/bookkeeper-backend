package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"bookkeeper-backend/models"
	"time"
)

func TestCreateGoal(t *testing.T) {
	goal := models.Goal{
		Name:        "Emergency Fund",
		HouseholdID: 1,
		TargetDate:  time.Now().AddDate(1, 0, 0),
		Category:    "Savings",
		Target:      5000,
		Progress:    1000,
		Notes:       "Save for emergencies",
	}
	body, _ := json.Marshal(goal)
	req := httptest.NewRequest("POST", "/goals/", bytes.NewReader(body))
	req = req.WithContext(mockUserContext(req.Context()))
	w := httptest.NewRecorder()
	createGoal(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Errorf("expected 200 or 201, got %d", w.Code)
	}
}