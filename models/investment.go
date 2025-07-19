package models

import (
	"time"
)

type Investment struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	AccountID   uint      `json:"accountId" gorm:"not null"`
	UserID      uint      `json:"userId" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null"`
	Institution string    `json:"institution" gorm:"not null"`
	Value       float64   `json:"value" gorm:"default:0"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}