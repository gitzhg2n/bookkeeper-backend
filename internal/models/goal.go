package models

import "time"

type Goal struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index"`
	Name      string    `gorm:"size:255"`
	TargetCents int64   `gorm:"not null"`
	CurrentCents int64  `gorm:"not null"`
	DueDate   *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
