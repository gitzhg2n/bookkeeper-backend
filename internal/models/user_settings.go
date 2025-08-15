package models

// NotificationPreferences defines user delivery channel preferences
type NotificationPreferences struct {
	InApp   bool `json:"in_app"`
	Email   bool `json:"email"`
	Push    bool `json:"push"`
}

type UserSettings struct {
	ID                        uint                   `gorm:"primaryKey"`
	UserID                    uint                   `gorm:"index;unique"`
	LargeTransactionThreshold int64                  `gorm:"default:10000"`
	NotificationPreferences   NotificationPreferences `gorm:"-" json:"notification_preferences,omitempty"`
}
