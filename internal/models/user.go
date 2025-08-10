package models

import (
	"time"
)

type User struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Email             string    `gorm:"uniqueIndex;size:255;not null" json:"email"`
	PasswordHash      []byte    `json:"-"`
	EncryptedDEK      []byte    `json:"-"`
	DEKSalt           []byte    `json:"-"`
	DEKNonce          []byte    `json:"-"`
	ArgonMemoryKiB    uint32    `json:"-"`
	ArgonTime         uint32    `json:"-"`
	ArgonParallelism  uint8     `json:"-"`
	ArgonSalt         []byte    `json:"-"`
	ArgonKeyLength    uint32    `json:"-"`
	KDFVersion        int       `json:"-"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	// Future: Recovery fields, household membership, roles, plan, etc.
}