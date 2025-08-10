package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestCoreFlow is a comprehensive integration test that tests the entire vertical slice:
// register → create household → create category → create account → upsert budget → create transaction → list budgets
// It verifies that actual_cents reflects the created transaction in the budget listing.
func TestCoreFlow(t *testing.T) {
	env := setupTest(t)

	// Step 1: Register a user
	regBody := `{"email":"test@example.com","password":"VerySecurePass1!"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/v1/auth/register", bytes.NewBufferString(regBody))
	env.Server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("register failed: expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	var authResp struct {
		Data struct {
			AccessToken string `json:"access_token"`
			UserID      uint   `json:"user_id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &authResp); err != nil {
		t.Fatalf("failed to unmarshal auth response: %v", err)
	}
	if authResp.Data.AccessToken == "" {
		t.Fatalf("missing access token")
	}
	
	accessToken := authResp.Data.AccessToken
	userID := authResp.Data.UserID
	authHeader := fmt.Sprintf("Bearer %s", accessToken)

	// Step 2: Create a household
	householdBody := `{"name":"Test Household"}`
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/v1/households", bytes.NewBufferString(householdBody))
	req.Header.Set("Authorization", authHeader)
	env.Server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("create household failed: expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	var householdResp struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &householdResp); err != nil {
		t.Fatalf("failed to unmarshal household response: %v", err)
	}
	householdID := householdResp.Data.ID

	// Step 3: Create a category
	categoryBody := `{"name":"Food"}`
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", fmt.Sprintf("/v1/households/%d/categories", householdID), bytes.NewBufferString(categoryBody))
	req.Header.Set("Authorization", authHeader)
	env.Server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("create category failed: expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	var categoryResp struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &categoryResp); err != nil {
		t.Fatalf("failed to unmarshal category response: %v", err)
	}
	categoryID := categoryResp.Data.ID

	// Step 4: Create an account
	accountBody := `{"name":"Checking Account","type":"checking","currency":"USD","opening_balance_cents":100000}`
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", fmt.Sprintf("/v1/households/%d/accounts", householdID), bytes.NewBufferString(accountBody))
	req.Header.Set("Authorization", authHeader)
	env.Server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("create account failed: expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	var accountResp struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &accountResp); err != nil {
		t.Fatalf("failed to unmarshal account response: %v", err)
	}
	accountID := accountResp.Data.ID

	// Step 5: Upsert a budget for current month
	currentMonth := time.Now().Format("2006-01")
	budgetBody := fmt.Sprintf(`{"month":"%s","category_id":%d,"planned_cents":50000}`, currentMonth, categoryID)
	w = httptest.NewRecorder()
	req = httptest.NewRequest("PUT", fmt.Sprintf("/v1/households/%d/budgets", householdID), bytes.NewBufferString(budgetBody))
	req.Header.Set("Authorization", authHeader)
	env.Server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("upsert budget failed: expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	// Step 6: Create a transaction
	transactionBody := fmt.Sprintf(`{"amount_cents":2500,"currency":"USD","category_id":%d,"memo":"Groceries","occurred_at":"%s"}`, 
		categoryID, time.Now().Format(time.RFC3339))
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", fmt.Sprintf("/v1/accounts/%d/transactions", accountID), bytes.NewBufferString(transactionBody))
	req.Header.Set("Authorization", authHeader)
	env.Server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("create transaction failed: expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	// Step 7: List budgets and verify actual_cents > 0
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", fmt.Sprintf("/v1/households/%d/budgets?month=%s", householdID, currentMonth), nil)
	req.Header.Set("Authorization", authHeader)
	env.Server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list budgets failed: expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	var budgetListResp struct {
		Data []struct {
			ID           uint   `json:"id"`
			Month        string `json:"month"`
			CategoryID   uint   `json:"category_id"`
			PlannedCents int64  `json:"planned_cents"`
			ActualCents  int64  `json:"actual_cents"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &budgetListResp); err != nil {
		t.Fatalf("failed to unmarshal budget list response: %v", err)
	}

	if len(budgetListResp.Data) == 0 {
		t.Fatalf("expected at least one budget, got none")
	}

	budget := budgetListResp.Data[0]
	if budget.ActualCents <= 0 {
		t.Fatalf("expected actual_cents > 0, got %d", budget.ActualCents)
	}

	// Verify the actual_cents matches our transaction amount
	if budget.ActualCents != 2500 {
		t.Fatalf("expected actual_cents to be 2500, got %d", budget.ActualCents)
	}

	// Verify other budget fields
	if budget.PlannedCents != 50000 {
		t.Fatalf("expected planned_cents to be 50000, got %d", budget.PlannedCents)
	}
	if budget.CategoryID != categoryID {
		t.Fatalf("expected category_id to be %d, got %d", categoryID, budget.CategoryID)
	}
	if budget.Month != currentMonth {
		t.Fatalf("expected month to be %s, got %s", currentMonth, budget.Month)
	}

	t.Logf("Core flow test completed successfully!")
	t.Logf("User ID: %d, Household ID: %d, Category ID: %d, Account ID: %d", userID, householdID, categoryID, accountID)
	t.Logf("Budget - Planned: %d cents, Actual: %d cents", budget.PlannedCents, budget.ActualCents)
}