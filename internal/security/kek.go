package security

import (
 copilot/fix-184f7982-e511-4e6f-9dc2-305d1c6b4c15
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

"crypto/hmac"
"crypto/sha256"
"encoding/hex"
"errors"
)

var (
ErrInvalidKey = errors.New("invalid key provided")
)

// DeriveKEK derives a 32-byte Key Encryption Key (KEK) from a password-derived key
// and a context string using HMAC-SHA256. This provides key separation between
// password verification material and DEK encryption material.
func DeriveKEK(passwordKey []byte, context string) ([]byte, error) {
if len(passwordKey) == 0 {
return nil, ErrInvalidKey
}
h := hmac.New(sha256.New, passwordKey)
h.Write([]byte(context))
return h.Sum(nil), nil // 32 bytes
}

// DeriveKEKHex returns the KEK as a hex-encoded string.
func DeriveKEKHex(passwordKey []byte, context string) (string, error) {
kek, err := DeriveKEK(passwordKey, context)
if err != nil {
return "", err
}
return hex.EncodeToString(kek), nil
 main
}