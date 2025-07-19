package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"bookkeeper-backend/models"
)

func TestCreateTransaction(t *testing.T) {
	tx := models.Transaction{
		AccountID: 1,
		Date:      models.Now(), // Assuming you have a helper for now, or use time.Now()
		Category:  "Groceries",
		Status:    "completed",
		Amount:    50.75,
		Notes:     "Weekly groceries",
	}
	body, _ := json.Marshal(tx)
	req := httptest.NewRequest("POST", "/transactions/", bytes.NewReader(body))
	req = req.WithContext(mockUserContext(req.Context()))
	w := httptest.NewRecorder()
	createTransaction(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Errorf("expected 200 or 201, got %d", w.Code)
	}
}

// Add similar tests for GET, PUT, DELETE as needed

func mockUserContext(ctx interface{}) interface{} {
	// Extend this to mock userID, accountIDs, etc.
	return ctx
}