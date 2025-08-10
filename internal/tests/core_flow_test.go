package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bookkeeper-backend/internal/models"
)

// TestCoreFlow tests the complete vertical slice: 
// register → create household → create category → create account → upsert budget → create transaction → list budgets
// Verifying that budget actual_cents reflects the created transaction
func TestCoreFlow(t *testing.T) {
	env := setupTest(t)

	// Step 1: Register user
	regBody := `{"email":"flow@example.com","password":"SecurePass123!"}`
	regResp := makeRequest(t, env, "POST", "/v1/auth/register", regBody)
	if regResp.Code != http.StatusOK {
		t.Fatalf("register failed: %d %s", regResp.Code, regResp.Body.String())
	}

	var authData struct {
		Data struct {
			AccessToken string `json:"access_token"`
			UserID      uint   `json:"user_id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(regResp.Body).Decode(&authData); err != nil {
		t.Fatalf("decode register response: %v", err)
	}
	accessToken := authData.Data.AccessToken
	userID := authData.Data.UserID

	// Step 2: Create household
	householdBody := `{"name":"Test Household"}`
	householdResp := makeAuthRequest(t, env, "POST", "/v1/households", householdBody, accessToken)
	if householdResp.Code != http.StatusOK {
		t.Fatalf("create household failed: %d %s", householdResp.Code, householdResp.Body.String())
	}

	var householdData struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(householdResp.Body).Decode(&householdData); err != nil {
		t.Fatalf("decode household response: %v", err)
	}
	householdID := householdData.Data.ID

	// Step 3: Create category
	categoryBody := `{"name":"Groceries"}`
	categoryResp := makeAuthRequest(t, env, "POST", "/v1/households/"+uintToString(householdID)+"/categories", categoryBody, accessToken)
	if categoryResp.Code != http.StatusOK {
		t.Fatalf("create category failed: %d %s", categoryResp.Code, categoryResp.Body.String())
	}

	var categoryData struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(categoryResp.Body).Decode(&categoryData); err != nil {
		t.Fatalf("decode category response: %v", err)
	}
	categoryID := categoryData.Data.ID

	// Step 4: Create account
	accountBody := `{"name":"Checking Account","account_type":"checking","balance_cents":100000}`
	accountResp := makeAuthRequest(t, env, "POST", "/v1/households/"+uintToString(householdID)+"/accounts", accountBody, accessToken)
	if accountResp.Code != http.StatusOK {
		t.Fatalf("create account failed: %d %s", accountResp.Code, accountResp.Body.String())
	}

	var accountData struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(accountResp.Body).Decode(&accountData); err != nil {
		t.Fatalf("decode account response: %v", err)
	}
	accountID := accountData.Data.ID

	// Step 5: Upsert budget (create a budget with planned amount)
	budgetBody := `{"month":"2024-01","category_id":` + uintToString(categoryID) + `,"planned_cents":50000}`
	budgetResp := makeAuthRequest(t, env, "PUT", "/v1/households/"+uintToString(householdID)+"/budgets", budgetBody, accessToken)
	if budgetResp.Code != http.StatusOK {
		t.Fatalf("upsert budget failed: %d %s", budgetResp.Code, budgetResp.Body.String())
	}

	// Step 6: Create transaction
	transactionBody := `{"amount_cents":2500,"memo":"Grocery shopping","occurred_at":"2024-01-15T10:00:00Z","category_id":` + uintToString(categoryID) + `}`
	transactionResp := makeAuthRequest(t, env, "POST", "/v1/accounts/"+uintToString(accountID)+"/transactions", transactionBody, accessToken)
	if transactionResp.Code != http.StatusOK {
		t.Fatalf("create transaction failed: %d %s", transactionResp.Code, transactionResp.Body.String())
	}

	// Step 7: List budgets and verify actual_cents > 0
	budgetListResp := makeAuthRequest(t, env, "GET", "/v1/households/"+uintToString(householdID)+"/budgets?month=2024-01", "", accessToken)
	if budgetListResp.Code != http.StatusOK {
		t.Fatalf("list budgets failed: %d %s", budgetListResp.Code, budgetListResp.Body.String())
	}

	var budgetListData struct {
		Data []struct {
			ID           uint  `json:"id"`
			CategoryID   uint  `json:"category_id"`
			PlannedCents int64 `json:"planned_cents"`
			ActualCents  int64 `json:"actual_cents"`
		} `json:"data"`
	}
	if err := json.NewDecoder(budgetListResp.Body).Decode(&budgetListData); err != nil {
		t.Fatalf("decode budget list response: %v", err)
	}

	// Find our budget and verify actual_cents reflects the transaction
	found := false
	for _, budget := range budgetListData.Data {
		if budget.CategoryID == categoryID {
			found = true
			if budget.ActualCents <= 0 {
				t.Errorf("Expected budget actual_cents > 0, got %d", budget.ActualCents)
			}
			if budget.PlannedCents != 50000 {
				t.Errorf("Expected planned_cents 50000, got %d", budget.PlannedCents)
			}
			// The actual_cents should be 2500 (the transaction amount)
			if budget.ActualCents != 2500 {
				t.Errorf("Expected actual_cents 2500, got %d", budget.ActualCents)
			}
			break
		}
	}
	if !found {
		t.Error("Budget not found in list")
	}

	// Bonus: Verify user is properly associated with household
	var userHousehold models.HouseholdMember
	if err := env.DB.Where("user_id = ? AND household_id = ?", userID, householdID).First(&userHousehold).Error; err != nil {
		t.Errorf("User should be associated with household: %v", err)
	}
}

// Helper functions for the core flow test
func makeRequest(t *testing.T, env *testEnv, method, path, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	env.Server.ServeHTTP(w, req)
	return w
}

func makeAuthRequest(t *testing.T, env *testEnv, method, path, body, token string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	env.Server.ServeHTTP(w, req)
	return w
}

func uintToString(u uint) string {
	return fmt.Sprintf("%d", u)
}