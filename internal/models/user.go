package models

import "time"

type User struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Email            string    `gorm:"uniqueIndex;size:255;not null" json:"email"`
	PasswordHash     []byte    `json:"-"`
	EncryptedDEK     []byte    `json:"-"`
	DEKNonce         []byte    `json:"-"`
	ArgonMemoryKiB   uint32    `json:"-"`
	ArgonTime        uint32    `json:"-"`
	ArgonParallelism uint8     `json:"-"`
	ArgonSalt        []byte    `json:"-"`
	ArgonKeyLength   uint32    `json:"-"`
	KDFVersion       int       `json:"-"`
	Plan             string    `gorm:"size:32;default:'free'" json:"plan"` // free, premium, selfhost
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type RefreshToken struct {
	ID           string  `gorm:"primaryKey;size:64"`
	UserID       uint    `gorm:"index"`
	User         User    `gorm:"constraint:OnDelete:CASCADE"`
	ExpiresAt    int64   `gorm:"index"`
	RevokedAt    *int64
	ReplacedByID *string
	CreatedAt    time.Time
}

type Household struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:255;not null"`
	CreatedBy uint      `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Members   []HouseholdMember
}

type HouseholdMember struct {
	ID          uint      `gorm:"primaryKey"`
	HouseholdID uint      `gorm:"index"`
	UserID      uint      `gorm:"index"`
	Role        string    `gorm:"size:32"`
	CreatedAt   time.Time
}

type Account struct {
	ID                  uint       `gorm:"primaryKey"`
	HouseholdID         uint       `gorm:"index"`
	Name                string     `gorm:"size:255"`
	Type                string     `gorm:"size:32"`
	Currency            string     `gorm:"size:8"`
	OpeningBalanceCents int64
	ArchivedAt          *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type Category struct {
	ID          uint      `gorm:"primaryKey"`
	HouseholdID uint      `gorm:"index"`
	Name        string    `gorm:"size:255"`
	ParentID    *uint     `gorm:"index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Budget struct {
	ID           uint      `gorm:"primaryKey"`
	HouseholdID  uint      `gorm:"index:uniq_budget,unique"`
	Month        string    `gorm:"size:7;index:uniq_budget,unique"`
	CategoryID   uint      `gorm:"index:uniq_budget,unique"`
	PlannedCents int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Transaction struct {
	ID          uint      `gorm:"primaryKey"`
	AccountID   uint      `gorm:"index"`
	UserID      *uint     `gorm:"index"`
	AmountCents int64
	Currency    string    `gorm:"size:8"`
	CategoryID  *uint     `gorm:"index"`
	Memo        string    `gorm:"size:1024"`
	OccurredAt  time.Time `gorm:"index"`
	CreatedAt   time.Time
}