package security

import (
	"crypto/hmac"
	"crypto/sha256"
)

// DeriveKEK derives a Key Encryption Key (KEK) from a password key using HMAC-SHA256.
// This separates the KEK from the password hash, providing better security architecture
// for encrypting/decrypting data encryption keys (DEKs).
//
// passwordKey: The key derived from password using Argon2 (from password.go)
// info: Context string for key derivation (e.g., "bookkeeper:dek:v1")
//
// Returns a 32-byte KEK suitable for encrypting/decrypting DEKs.
func DeriveKEK(passwordKey []byte, info string) []byte {
	// Use HMAC-SHA256 as a KDF to derive KEK from password key
	// This is similar to HKDF-Expand with a single step
	h := hmac.New(sha256.New, passwordKey)
	h.Write([]byte(info))
	return h.Sum(nil) // Returns 32 bytes (SHA256 output)
}