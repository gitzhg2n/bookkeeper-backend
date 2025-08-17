package models

import "time"

// CompoundAlertCondition allows multiple conditions to be combined with AND/OR
// Example: (A > 100 AND B < 50) OR (C == 200)
type CompoundAlertCondition struct {
	Operator   string                   `json:"operator"` // AND, OR
	Conditions []InvestmentAlertCondition `json:"conditions"`
}

// InvestmentAlertCondition represents a single alert condition (for compound alerts)
type InvestmentAlertCondition struct {
	AssetSymbol string  `json:"asset_symbol"`
	AlertType   string  `json:"alert_type"`
	Direction   string  `json:"direction"`
	Threshold   float64 `json:"threshold"`
}

// TimeWindow defines a period for time-based alert triggers
type TimeWindow struct {
	Duration string `json:"duration"` // e.g. "24h", "7d"
}

type InvestmentAlert struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"index"`
	AssetSymbol  string    `gorm:"size:32"`
	AlertType    string    `gorm:"size:32"` // price, percent_change, value
	Direction    string    `gorm:"size:8"`  // up, down
	Threshold    float64   // price or percent
	Active       bool      `gorm:"default:true"`
	Compound     *CompoundAlertCondition `gorm:"-" json:"compound,omitempty"` // not stored in DB directly
	TimeWindow   *TimeWindow             `gorm:"-" json:"time_window,omitempty"`
	CooldownMinutes int                  `gorm:"-" json:"cooldown_minutes,omitempty"`
	// Rule is the stored rule expression for the alert (persisted in DB)
	Rule          string                  `gorm:"size:1024" json:"rule,omitempty"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
