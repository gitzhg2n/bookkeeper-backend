package models

import (
	"context"
)

// GetHouseholdByID fetches household by ID with context
func GetHouseholdByID(ctx context.Context, id uint) *Household {
	var h Household
	if err := DB.WithContext(ctx).First(&h, id).Error; err != nil {
		return nil
	}
	return &h
}

// GetAccountByID fetches account by ID with context
func GetAccountByID(ctx context.Context, id uint) *Account {
	var a Account
	if err := DB.WithContext(ctx).First(&a, id).Error; err != nil {
		return nil
	}
	return &a
}

// GetBudgetByID fetches budget by ID with context
func GetBudgetByID(ctx context.Context, id uint) *Budget {
	var b Budget
	if err := DB.WithContext(ctx).First(&b, id).Error; err != nil {
		return nil
	}
	return &b
}

// GetGoalByID fetches goal by ID with context
func GetGoalByID(ctx context.Context, id uint) *Goal {
	var g Goal
	if err := DB.WithContext(ctx).First(&g, id).Error; err != nil {
		return nil
	}
	return &g
}

// GetInvestmentByID fetches investment by ID with context
func GetInvestmentByID(ctx context.Context, id uint) *Investment {
	var i Investment
	if err := DB.WithContext(ctx).First(&i, id).Error; err != nil {
		return nil
	}
	return &i
}

// GetTransactionByID fetches transaction by ID with context
func GetTransactionByID(ctx context.Context, id uint) *Transaction {
	var t Transaction
	if err := DB.WithContext(ctx).First(&t, id).Error; err != nil {
		return nil
	}
	return &t
}