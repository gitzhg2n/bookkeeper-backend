package models

import "time"

type Bill struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index"`
	Name      string    `gorm:"size:255"`
	AmountCents int64   `gorm:"not null"`
	DueDay    int       `gorm:"not null"` // Day of month
	NextDue   time.Time `gorm:"not null"`
	Recurring bool      `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
