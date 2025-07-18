package models

import (
	"time"
)

type Investment struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	AccountID   uint      `json:"accountId" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null"`
	Institution string    `json:"institution" gorm:"not null"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}