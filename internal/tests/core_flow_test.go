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

// TestCoreFlow tests the vertical slice: register → create household → create category → create account
// → upsert budget → create transaction → list budgets, and verifies actual_cents aggregation.
func TestCoreFlow(t *testing.T) {
	env := setupTest(t)

	// Register
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
		t.Fatalf("decode register: %v", err)
	}
	token := authData.Data.AccessToken
	userID := authData.Data.UserID

	// Household
	houseResp := makeAuthRequest(t, env, "POST", "/v1/households", `{"name":"Test Household"}`, token)
	if houseResp.Code != http.StatusOK {
		t.Fatalf("household create failed: %d %s", houseResp.Code, houseResp.Body.String())
	}
	var hData struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	_ = json.NewDecoder(houseResp.Body).Decode(&hData)
	householdID := hData.Data.ID

	// Category
	catResp := makeAuthRequest(t, env, "POST", fmt.Sprintf("/v1/households/%d/categories", householdID), `{"name":"Groceries"}`, token)
	if catResp.Code != http.StatusOK {
		t.Fatalf("category create failed: %d %s", catResp.Code, catResp.Body.String())
	}
	var cData struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	_ = json.NewDecoder(catResp.Body).Decode(&cData)
	categoryID := cData.Data.ID

	// Account
	accResp := makeAuthRequest(t, env, "POST", fmt.Sprintf("/v1/households/%d/accounts", householdID),
		`{"name":"Checking Account","type":"checking","currency":"USD","opening_balance_cents":100000}`, token)
	if accResp.Code != http.StatusOK {
		t.Fatalf("account create failed: %d %s", accResp.Code, accResp.Body.String())
	}
	var aData struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}
	_ = json.NewDecoder(accResp.Body).Decode(&aData)
	accountID := aData.Data.ID

	// Budget
	budgetResp := makeAuthRequest(t, env, "PUT", fmt.Sprintf("/v1/households/%d/budgets", householdID),
		fmt.Sprintf(`{"month":"2024-01","category_id":%d,"planned_cents":50000}`, categoryID), token)
	if budgetResp.Code != http.StatusOK {
		t.Fatalf("budget upsert failed: %d %s", budgetResp.Code, budgetResp.Body.String())
	}

	// Transaction
	trxResp := makeAuthRequest(t, env, "POST", fmt.Sprintf("/v1/accounts/%d/transactions", accountID),
		fmt.Sprintf(`{"amount_cents":2500,"currency":"USD","category_id":%d,"memo":"Grocery shopping","occurred_at":"2024-01-15T10:00:00Z"}`, categoryID),
		token)
	if trxResp.Code != http.StatusOK {
		t.Fatalf("transaction create failed: %d %s", trxResp.Code, trxResp.Body.String())
	}

	// Budgets list
	listResp := makeAuthRequest(t, env, "GET", fmt.Sprintf("/v1/households/%d/budgets?month=2024-01", householdID), "", token)
	if listResp.Code != http.StatusOK {
		t.Fatalf("budget list failed: %d %s", listResp.Code, listResp.Body.String())
	}
	var listData struct {
		Data []struct {
			CategoryID   uint  `json:"category_id"`
			PlannedCents int64 `json:"planned_cents"`
			ActualCents  int64 `json:"actual_cents"`
		} `json:"data"`
	}
	_ = json.NewDecoder(listResp.Body).Decode(&listData)
	found := false
	for _, b := range listData.Data {
		if b.CategoryID == categoryID {
			found = true
			if b.PlannedCents != 50000 {
				t.Errorf("expected planned 50000 got %d", b.PlannedCents)
			}
			if b.ActualCents != 2500 {
				t.Errorf("expected actual 2500 got %d", b.ActualCents)
			}
		}
	}
	if !found {
		t.Errorf("budget not found for category")
	}

	// Membership assert
	var hm models.HouseholdMember
	if err := env.DB.Where("user_id = ? AND household_id = ?", userID, householdID).First(&hm).Error; err != nil {
		t.Errorf("user should be member: %v", err)
	}
}

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