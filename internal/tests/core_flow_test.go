package tests

import (
"bytes
"encoding/jso
" "f
t" "net/h
tp" "net/http/http
est" "te
ting"
// TestCoreFlow is a comprehensive integration test for the full vertical slice:
// register → create household → create category → create account → upsert budget → create transaction → list budgets
// It verifies that actual_cents reflects the created expense transaction in the budget listing.
func TestCoreFlow(t *testing.T) 
  
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
if authResp.Data.AccessToken "" {
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
if householdID  0 {
	t.Fatalf("missing household ID")
}

// Step 3: Create a category
categoryBody := `{"name":"Groceries","type":"expense"}`
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
if categoryID 0 {
	t.Fatalf("missing category ID")
}

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
if accountID == 0 {
	t.Fatalf("missing account ID")
}

// Step 5: Upsert a budget for the current month
currentMonth := time.Now().Format("2006-01")
budgetUpsertBody := fmt.Sprintf(`{"month":"%s","category_id":%d,"planned_cents":50000}`, currentMonth, categoryID)
w = httptest.NewRecorder()
req = httptest.NewRequest("PUT", fmt.Sprintf("/v1/households/%d/budgets", householdID), bytes.NewBufferString(budgetUpsertBody))
req.Header.Set("Authorization", authHeader)
env.Server.ServeHTTP(w, req)
if w.Code != http.StatusOK {
	t.Fatalf("upsert budget failed: expected 200 got %d body=%s", w.Code, w.Body.String())
}

var budgetUpsertResp struct {
	Data struct {
		ID uint `json:"id"`
	} `json:"data"`
}
if err := json.Unmarshal(w.Body.Bytes(), &budgetUpsertResp); err != nil {
	t.Fatalf("failed to unmarshal budget upsert response: %v", err)
}
budgetID := budgetUpsertResp.Data.ID
if budgetID == 0 {
	t.Fatalf("missing budget ID")
}

// Step 6: Create an expense transaction (negative amount)
occurredAt := time.Now().Format(time.RFC3339)
transactionBody := fmt.Sprintf(`{"amount_cents":-2500,"currency":"USD","category_id":%d,"memo":"Grocery shopping","occurred_at":"%s"}`, categoryID, occurredAt)
w = httptest.NewRecorder()
req = httptest.NewRequest("POST", fmt.Sprintf("/v1/accounts/%d/transactions", accountID), bytes.NewBufferString(transactionBody))
req.Header.Set("Authorization", authHeader)
env.Server.ServeHTTP(w, req)
if w.Code != http.StatusOK {
	t.Fatalf("create transaction failed: expected 200 got %d body=%s", w.Code, w.Body.String())
}

// Step 7: List budgets and verify actual_cents reflects the expense
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
	t.Fatalf("no budgets returned")
}

// Find our budget and validate fields
found := false
for _, b := range budgetListResp.Data {
	if b.ID == budgetID {
		found = true
		if b.Month != currentMonth {
			t.Fatalf("expected month %s, got %s", currentMonth, b.Month)
		}
		if b.CategoryID != categoryID {
			t.Fatalf("expected category_id %d, got %d", categoryID, b.CategoryID)
		}
		if b.PlannedCents != 50000 {
			t.Fatalf("expected planned_cents 50000, got %d", b.PlannedCents)
		}
		// Expense should decrease actual_cents
		if b.ActualCents != -2500 {
			t.Fatalf("expected actual_cents -2500, got %d", b.ActualCents)
		}
		break
	}
}
if !found {
	t.Fatalf("created budget with ID %d not found in listing", budgetID)
}

t.Logf("Core flow test completed successfully! User ID: %d, Household ID: %d, Category ID: %d, Account ID: %d", userID, householdID, categoryID, accountID)
