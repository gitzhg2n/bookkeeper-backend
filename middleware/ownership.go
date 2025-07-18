package middleware

import (
	"bookkeeper-backend-go/models"
	"context"
)

func CheckHouseholdOwnership(ctx context.Context, userID uint, householdID uint) bool {
	household := models.GetHouseholdByID(ctx, householdID)
	return household != nil && household.OwnerID == userID
}

func CheckAccountOwnership(ctx context.Context, userID uint, accountID uint) bool {
	account := models.GetAccountByID(ctx, accountID)
	if account == nil {
		return false
	}
	return CheckHouseholdOwnership(ctx, userID, account.HouseholdID)
}

func CheckBudgetOwnership(ctx context.Context, userID uint, budgetID uint) bool {
	budget := models.GetBudgetByID(ctx, budgetID)
	if budget == nil {
		return false
	}
	return CheckHouseholdOwnership(ctx, userID, budget.HouseholdID)
}

func CheckGoalOwnership(ctx context.Context, userID uint, goalID uint) bool {
	goal := models.GetGoalByID(ctx, goalID)
	if goal == nil {
		return false
	}
	return CheckHouseholdOwnership(ctx, userID, goal.HouseholdID)
}

func CheckInvestmentOwnership(ctx context.Context, userID uint, investmentID uint) bool {
	investment := models.GetInvestmentByID(ctx, investmentID)
	if investment == nil {
		return false
	}
	return CheckAccountOwnership(ctx, userID, investment.AccountID)
}

func CheckTransactionOwnership(ctx context.Context, userID uint, transactionID uint) bool {
	transaction := models.GetTransactionByID(ctx, transactionID)
	if transaction == nil {
		return false
	}
	return CheckAccountOwnership(ctx, userID, transaction.AccountID)
}