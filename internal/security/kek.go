package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

var (
	ErrInvalidKey = errors.New("invalid key provided")
)

// DeriveKEK derives a Key Encryption Key (KEK) from a password key and context string
// using HMAC-SHA256 expansion. This provides proper key separation between the 
// password verification key and the data encryption key.
func DeriveKEK(passwordKey []byte, context string) ([]byte, error) {
	if len(passwordKey) == 0 {
		return nil, ErrInvalidKey
	}
	
	// Use HMAC-SHA256 to derive KEK from password key with context
	h := hmac.New(sha256.New, passwordKey)
	h.Write([]byte(context))
	kek := h.Sum(nil)
	
	return kek, nil
}

// DeriveKEKHex is a convenience function that returns the KEK as a hex string
func DeriveKEKHex(passwordKey []byte, context string) (string, error) {
	kek, err := DeriveKEK(passwordKey, context)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(kek), nil
}