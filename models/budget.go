package models

import (
	"time"
)

type Budget struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	HouseholdID uint      `json:"householdId" gorm:"not null"`
	UserID      uint      `json:"userId" gorm:"not null"`
	Amount      float64   `json:"amount" gorm:"not null"`
	Period      string    `json:"period" gorm:"not null"`
	Category    string    `json:"category" gorm:"not null"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}