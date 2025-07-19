package routes

import (
	"testing"
	"bookkeeper-backend/models"
)

func TestUpdateAccount(t *testing.T) {
	acc := models.Account{
		Name:        "Old Name",
		Type:        "checking",
		HouseholdID: 1,
		Institution: "Bank",
		Balance:     100.0,
	}
	models.DB.Create(&acc)
	reqPayload := AccountRequest{
		Name:    "New Name",
		Type:    "checking",
		Balance: 200.0,
	}
	body, _ := json.Marshal(reqPayload)
	req := httptest.NewRequest("PUT", "/accounts/"+strconv.Itoa(int(acc.ID)), bytes.NewReader(body))
	req = req.WithContext(mockUserContext(req.Context()))
	w := httptest.NewRecorder()
	updateAccount(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestDeleteAccount(t *testing.T) {
	acc := models.Account{
		Name:        "Delete Me",
		Type:        "checking",
		HouseholdID: 1,
		Institution: "Bank",
		Balance:     100.0,
	}
	models.DB.Create(&acc)
	req := httptest.NewRequest("DELETE", "/accounts/"+strconv.Itoa(int(acc.ID)), nil)
	req = req.WithContext(mockUserContext(req.Context()))
	w := httptest.NewRecorder()
	deleteAccount(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}