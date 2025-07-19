package models

import (
	"time"
)

type Account struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null"`
	HouseholdID uint      `json:"householdId" gorm:"not null"`
	Institution string    `json:"institution" gorm:"not null"`
	Balance     float64   `json:"balance" gorm:"default:0"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}