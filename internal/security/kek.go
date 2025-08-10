package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

var ErrInvalidKey = errors.New("invalid key provided")

// DeriveKEK derives a 32-byte KEK from a password-derived key + context using HMAC-SHA256.
// This provides key separation so the password-derived key can be used to derive
// multiple independent subkeys (e.g., DEK wrap key, future verifier material).
func DeriveKEK(passwordKey []byte, context string) ([]byte, error) {
	if len(passwordKey) == 0 {
		return nil, ErrInvalidKey
	}
	h := hmac.New(sha256.New, passwordKey)
	h.Write([]byte(context))
	return h.Sum(nil), nil
}

// DeriveKEKHex helper for logging/debug (avoid storing this; shown here for completeness).
func DeriveKEKHex(passwordKey []byte, context string) (string, error) {
	kek, err := DeriveKEK(passwordKey, context)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(kek), nil
}