package models

import (
	"context"
)

// Example: Helper functions for ownership checks
func GetHouseholdByID(ctx context.Context, id uint) *Household {
	var h Household
	if err := DB.First(&h, id).Error; err != nil {
		return nil
	}
	return &h
}
func GetAccountByID(ctx context.Context, id uint) *Account {
	var a Account
	if err := DB.First(&a, id).Error; err != nil {
		return nil
	}
	return &a
}
func GetBudgetByID(ctx context.Context, id uint) *Budget {
	var b Budget
	if err := DB.First(&b, id).Error; err != nil {
		return nil
	}
	return &b
}
func GetGoalByID(ctx context.Context, id uint) *Goal {
	var g Goal
	if err := DB.First(&g, id).Error; err != nil {
		return nil
	}
	return &g
}
func GetInvestmentByID(ctx context.Context, id uint) *Investment {
	var i Investment
	if err := DB.First(&i, id).Error; err != nil {
		return nil
	}
	return &i
}
func GetTransactionByID(ctx context.Context, id uint) *Transaction {
	var t Transaction
	if err := DB.First(&t, id).Error; err != nil {
		return nil
	}
	return &t
}