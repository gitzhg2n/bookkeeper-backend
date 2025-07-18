package models

import (
	"time"
)

type HouseholdMember struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	HouseholdID uint      `json:"householdId" gorm:"not null"`
	CreatedAt   time.Time `json:"createdAt"`
}