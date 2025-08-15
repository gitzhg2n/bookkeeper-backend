package models

import "time"

type Transaction struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	AccountID   uint      `json:"account_id"`
	UserID      *uint     `json:"user_id"`
	AmountCents int64     `json:"amount_cents"`
	Currency    string    `json:"currency"`
	CategoryID  *uint     `json:"category_id"`
	Memo        string    `json:"memo"`
	OccurredAt  time.Time `json:"occurred_at"`
}
