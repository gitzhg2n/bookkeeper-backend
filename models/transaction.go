package models

import (
	"time"
)

type Transaction struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	AccountID uint      `json:"accountId" gorm:"not null"`
	Date      time.Time `json:"date" gorm:"not null"`
	Category  string    `json:"category" gorm:"not null"`
	Status    string    `json:"status" gorm:"not null"`
	Amount    float64   `json:"amount" gorm:"not null"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}