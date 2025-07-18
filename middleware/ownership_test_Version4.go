package middleware

import (
	"context"
	"testing"
	"bookkeeper-backend-go/models"
)

// Mock model functions for testing
func mockGetHouseholdByID(ctx context.Context, id uint) *models.Household {
	if id == 1 {
		return &models.Household{ID: 1, OwnerID: 42}
	}
	return nil
}
func mockGetAccountByID(ctx context.Context, id uint) *models.Account {
	if id == 2 {
		return &models.Account{ID: 2, HouseholdID: 1}
	}
	return nil
}
func mockGetBudgetByID(ctx context.Context, id uint) *models.Budget {
	if id == 3 {
		return &models.Budget{ID: 3, HouseholdID: 1}
	}
	return nil
}
func mockGetGoalByID(ctx context.Context, id uint) *models.Goal {
	if id == 4 {
		return &models.Goal{ID: 4, HouseholdID: 1}
	}
	return nil
}
func mockGetInvestmentByID(ctx context.Context, id uint) *models.Investment {
	if id == 5 {
		return &models.Investment{ID: 5, AccountID: 2}
	}
	return nil
}
func mockGetTransactionByID(ctx context.Context, id uint) *models.Transaction {
	if id == 6 {
		return &models.Transaction{ID: 6, AccountID: 2}
	}
	return nil
}

func TestOwnershipChecks(t *testing.T) {
	models.GetHouseholdByID = mockGetHouseholdByID
	models.GetAccountByID = mockGetAccountByID
	models.GetBudgetByID = mockGetBudgetByID
	models.GetGoalByID = mockGetGoalByID
	models.GetInvestmentByID = mockGetInvestmentByID
	models.GetTransactionByID = mockGetTransactionByID
	ctx := context.Background()

	if !CheckHouseholdOwnership(ctx, 42, 1) {
		t.Error("Expected household ownership to be true")
	}
	if CheckHouseholdOwnership(ctx, 43, 1) {
		t.Error("Expected household ownership to be false")
	}
	if !CheckAccountOwnership(ctx, 42, 2) {
		t.Error("Expected account ownership to be true")
	}
	if !CheckBudgetOwnership(ctx, 42, 3) {
		t.Error("Expected budget ownership to be true")
	}
	if !CheckGoalOwnership(ctx, 42, 4) {
		t.Error("Expected goal ownership to be true")
	}
	if !CheckInvestmentOwnership(ctx, 42, 5) {
		t.Error("Expected investment ownership to be true")
	}
	if !CheckTransactionOwnership(ctx, 42, 6) {
		t.Error("Expected transaction ownership to be true")
	}
}