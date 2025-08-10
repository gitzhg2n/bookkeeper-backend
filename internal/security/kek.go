package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// DeriveKEK derives a Key Encryption Key (KEK) from a password key using HMAC-SHA256
// This follows the pattern of HMAC-based key derivation for generating secondary keys
func DeriveKEK(passwordKey []byte, context string) []byte {
	h := hmac.New(sha256.New, passwordKey)
	h.Write([]byte(context))
	return h.Sum(nil)
}

// DeriveKEKHex is a convenience function that returns the KEK as a hex string
func DeriveKEKHex(passwordKey []byte, context string) string {
	kek := DeriveKEK(passwordKey, context)
	return hex.EncodeToString(kek)
}