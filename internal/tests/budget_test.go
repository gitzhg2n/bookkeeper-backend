package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBudgetLifecycle(t *testing.T) {
	env := setupTest(t)

	// Register user
	reg := `{"email":"b@example.com","password":"StrongPassw0rd!"}`
	w := httptest.NewRecorder()
	env.Server.ServeHTTP(w, httptest.NewRequest("POST", "/v1/auth/register", bytes.NewBufferString(reg)))
	if w.Code != http.StatusOK {
		t.Fatalf("register failed: %d %s", w.Code, w.Body.String())
	}
	var regResp struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &regResp)
	token := regResp.Data.AccessToken
	if token == "" {
		t.Fatal("no token")
	}

	// Create household
	w = httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/v1/households", bytes.NewBufferString(`{"name":"Home"}`))
	req.Header.Set("Authorization", "Bearer "+token)
	env.Server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("household create failed: %d %s", w.Code, w.Body.String())
	}
	var houseResp struct {
		Data struct {
			ID int `json:"id"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &houseResp)
	hID := houseResp.Data.ID

	// Create category
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/v1/households/"+fmt.Sprint(hID)+"/categories", bytes.NewBufferString(`{"name":"Groceries"}`))
	req.Header.Set("Authorization", "Bearer "+token)
	env.Server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("category create failed: %d %s", w.Code, w.Body.String())
	}
	var catResp struct {
		Data struct {
			ID int `json:"id"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &catResp)
	catID := catResp.Data.ID

	// Upsert budget
	w = httptest.NewRecorder()
	req = httptest.NewRequest("PUT", "/v1/households/"+fmt.Sprint(hID)+"/budgets",
		bytes.NewBufferString(fmt.Sprintf(`{"month":"2025-08","category_id":%d,"planned_cents":50000}`, catID)))
	req.Header.Set("Authorization", "Bearer "+token)
	env.Server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("budget upsert failed: %d %s", w.Code, w.Body.String())
	}

	// List budgets
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/v1/households/"+fmt.Sprint(hID)+"/budgets?month=2025-08", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	env.Server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("budget list failed: %d %s", w.Code, w.Body.String())
	}
}