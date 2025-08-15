package models

import "time"

// AlertHistory records when an alert was triggered and why
type AlertHistory struct {
	ID          uint      `gorm:"primaryKey"`
	AlertID     uint      `gorm:"index"`
	UserID      uint      `gorm:"index"`
	TriggeredAt time.Time
	Details     string    `gorm:"size:1024"`
}
