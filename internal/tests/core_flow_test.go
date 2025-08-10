package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCoreFlow tests the entire vertical slice:
// register → create household → create category → create account → upsert budget → create transaction → list budgets verifying actual_cents reflects transaction
func TestCoreFlow(t *testing.T) {
	env := setupTest(t)

	// Step 1: Register a user
	regBody := `{"email":"user@example.com","password":"VerySecurePass1!"}`
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
		t.Fatalf("failed to parse register response: %v", err)
	}
	accessToken := authResp.Data.AccessToken
	if accessToken == "" {
		t.Fatalf("missing access token")
	}

	// Step 2: Create household
	hhBody := `{"name":"Test Household"}`
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/v1/households", bytes.NewBufferString(hhBody))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	env.Server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("create household failed: expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	var hhResp struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &hhResp); err != nil {
		t.Fatalf("failed to parse household response: %v", err)
	}
	householdID := hhResp.Data.ID
	if householdID == 0 {
		t.Fatalf("missing household ID")
	}

	// Step 3: Create category
	catBody := `{"name":"Groceries","type":"expense"}`
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", fmt.Sprintf("/v1/households/%d/categories", householdID), bytes.NewBufferString(catBody))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	env.Server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("create category failed: expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	var catResp struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &catResp); err != nil {
		t.Fatalf("failed to parse category response: %v", err)
	}
	categoryID := catResp.Data.ID
	if categoryID == 0 {
		t.Fatalf("missing category ID")
	}

	// Step 4: Create account
	accBody := `{"name":"Checking Account","type":"checking"}`
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", fmt.Sprintf("/v1/households/%d/accounts", householdID), bytes.NewBufferString(accBody))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	env.Server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("create account failed: expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	var accResp struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &accResp); err != nil {
		t.Fatalf("failed to parse account response: %v", err)
	}
	accountID := accResp.Data.ID
	if accountID == 0 {
		t.Fatalf("missing account ID")
	}

	// Step 5: Upsert budget (create budget for the category)
	budgetBody := fmt.Sprintf(`{"month":"2024-12","category_id":%d,"planned_cents":50000}`, categoryID)
	w = httptest.NewRecorder()
	req = httptest.NewRequest("PUT", fmt.Sprintf("/v1/households/%d/budgets", householdID), bytes.NewBufferString(budgetBody))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	env.Server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("upsert budget failed: expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	var budgetResp struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &budgetResp); err != nil {
		t.Fatalf("failed to parse budget response: %v", err)
	}
	budgetID := budgetResp.Data.ID
	if budgetID == 0 {
		t.Fatalf("missing budget ID")
	}

	// Step 6: Create transaction (this should affect budget actual_cents)
	txnBody := fmt.Sprintf(`{"amount_cents":-2500,"memo":"Grocery shopping","category_id":%d,"occurred_at":"2024-12-15T12:00:00Z"}`, categoryID)
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", fmt.Sprintf("/v1/accounts/%d/transactions", accountID), bytes.NewBufferString(txnBody))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	env.Server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("create transaction failed: expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	// Step 7: List budgets and verify actual_cents > 0
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", fmt.Sprintf("/v1/households/%d/budgets?month=2024-12", householdID), nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	env.Server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list budgets failed: expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	var budgetsResp struct {
		Data []struct {
			ID           uint  `json:"id"`
			ActualCents  int64 `json:"actual_cents"`
			PlannedCents int64 `json:"planned_cents"`
			CategoryID   uint  `json:"category_id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &budgetsResp); err != nil {
		t.Fatalf("failed to parse budgets response: %v", err)
	}

	if len(budgetsResp.Data) == 0 {
		t.Fatalf("no budgets returned")
	}

	// Find our budget and verify actual_cents reflects transaction
	found := false
	for _, budget := range budgetsResp.Data {
		if budget.ID == budgetID {
			found = true
			// The actual_cents should be -2500 (reflecting the expense transaction)
			if budget.ActualCents != -2500 {
				t.Fatalf("budget actual_cents should be -2500 (expense), got %d", budget.ActualCents)
			}
			if budget.PlannedCents != 50000 {
				t.Fatalf("budget planned_cents should be 50000, got %d", budget.PlannedCents)
			}
			if budget.CategoryID != categoryID {
				t.Fatalf("budget category_id should be %d, got %d", categoryID, budget.CategoryID)
			}
			t.Logf("Success: Budget actual_cents is %d, correctly reflecting the expense transaction", budget.ActualCents)
			break
		}
	}

	if !found {
		t.Fatalf("created budget with ID %d not found in list", budgetID)
	}

	t.Logf("Core flow test passed! Budget actual_cents correctly reflects transaction amount.")
}