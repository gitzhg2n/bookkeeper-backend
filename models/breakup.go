package models

import (
	"time"
)

type BreakupRequest struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Reason        string    `json:"reason" gorm:"not null"`
	Account       string    `json:"account" gorm:"not null"`
	TransferAmount float64  `json:"transferAmount" gorm:"not null"`
	CreatedAt     time.Time `json:"createdAt"`
}