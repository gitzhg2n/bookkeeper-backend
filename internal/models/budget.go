package models

type Budget struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	HouseholdID uint   `json:"household_id"`
	Month       string `json:"month"`
	CategoryID  uint   `json:"category_id"`
	PlannedCents int64 `json:"planned_cents"`
}
