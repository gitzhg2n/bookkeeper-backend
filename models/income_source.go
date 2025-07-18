package models

import (
	"time"
)

type IncomeSource struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null"` // "W-2", "1099", "K-1", "Other"
	Amount      float64   `json:"amount" gorm:"not null"`
	HouseholdID uint      `json:"householdId" gorm:"not null"`
	Frequency   string    `json:"frequency" gorm:"not null;default:monthly"` // "monthly", "annual"
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}