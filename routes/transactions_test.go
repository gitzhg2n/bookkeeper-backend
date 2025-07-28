package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	
	"bookkeeper-backend/models"
)

func TestCreateTransaction(t *testing.T) {
	tx := models.Transaction{
		AccountID: 1,
		Date:      time.Now(),
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
