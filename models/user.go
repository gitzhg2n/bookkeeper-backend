package models

import (
	"time"
)

type User struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Email            string    `json:"email" gorm:"unique;not null"`
	PasswordHash     string    `json:"passwordHash" gorm:"not null"`
	RecoverySeedHash string    `json:"recoverySeedHash" gorm:"not null"`
	Role             string    `json:"role" gorm:"not null;default:user"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}