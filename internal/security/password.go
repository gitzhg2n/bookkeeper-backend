package security

import (
	"crypto/subtle"
	"errors"

	"golang.org/x/crypto/argon2"
)

type ArgonParams struct {
	MemoryKiB    uint32
	Time         uint32
	Parallelism  uint8
	SaltLength   uint32
	KeyLength    uint32
}

func DeriveKey(password string, salt []byte, p ArgonParams) []byte {
	return argon2.IDKey([]byte(password), salt, p.Time, p.MemoryKiB, p.Parallelism, p.KeyLength)
}

func ConstantTimeCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare(a, b) == 1
}

var (
	ErrWeakPassword = errors.New("password does not meet minimum strength requirements")
)

func ValidatePasswordStrength(password string, allowInsecure bool) error {
	if allowInsecure {
		return nil
	}
	// Basic initial rule set (can expand later)
	if len(password) < 12 {
		return ErrWeakPassword
	}
	return nil
}