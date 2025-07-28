package models

import (
	"time"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Email            string    `json:"email" gorm:"unique;not null"`
	PasswordHash     string    `json:"-" gorm:"not null"`
	RecoverySeedHash string    `json:"-" gorm:"not null"`
	Role             string    `json:"role" gorm:"not null;default:user"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// HashRecoverySeed hashes a recovery seed using bcrypt
func HashRecoverySeed(seed string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(seed), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckRecoverySeed verifies a recovery seed against its hash
func CheckRecoverySeed(seed, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(seed))
	return err == nil
}