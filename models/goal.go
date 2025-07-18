package models

import (
	"time"
)

type Goal struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	HouseholdID uint      `json:"householdId" gorm:"not null"`
	TargetDate  time.Time `json:"targetDate" gorm:"not null"`
	Category    string    `json:"category" gorm:"not null"`
	Target      float64   `json:"target"`
	Progress    float64   `json:"progress"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}