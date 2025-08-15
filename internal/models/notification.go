package models

import (
	"time"
)

// NotificationType defines the type of notification (budget, transaction, goal, etc.)
type NotificationType string

const (
	NotificationTypeBudget            NotificationType = "budget"
	NotificationTypeTransaction       NotificationType = "transaction"
	NotificationTypeGoal              NotificationType = "goal"
	NotificationTypeInvestmentAlert   NotificationType = "investment_alert"
)

// Notification represents a user notification/alert
// All sensitive info should be in the message, not in type or metadata
// Only the user who owns the notification can access it
//
type Notification struct {
	ID        int64            `json:"id" db:"id"`
	UserID    int64            `json:"user_id" db:"user_id"`
	Type      NotificationType `json:"type" db:"type"`
	Title     string           `json:"title" db:"title"`
	Message   string           `json:"message" db:"message"`
	Read      bool             `json:"read" db:"read"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
}
