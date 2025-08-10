package security

import (
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/chacha20poly1305"
)

const (
	DEKLength = 32
	SaltLengthDefault = 16
)

type EncryptedDEK struct {
	Ciphertext []byte
	Salt       []byte
	Nonce      []byte
}

func RandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

// WrapDEK encrypts a freshly generated DEK with the KEK (derived from password)
func WrapDEK(kek []byte) (dek []byte, out EncryptedDEK, err error) {
	dek, err = RandomBytes(DEKLength)
	if err != nil {
		return nil, EncryptedDEK{}, err
	}
	aead, err := chacha20poly1305.NewX(kek)
	if err != nil {
		return nil, EncryptedDEK{}, err
	}
	nonce, err := RandomBytes(chacha20poly1305.NonceSizeX)
	if err != nil {
		return nil, EncryptedDEK{}, err
	}
	ct := aead.Seal(nil, nonce, dek, nil)
	out = EncryptedDEK{
		Ciphertext: ct,
		Nonce:      nonce,
		// Note: salt is associated with deriving kek; stored separately outside here.
	}
	return dek, out, nil
}

func UnwrapDEK(kek []byte, enc EncryptedDEK) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(kek)
	if err != nil {
		return nil, err
	}
	plain, err := aead.Open(nil, enc.Nonce, enc.Ciphertext, nil)
	if err != nil {
		return nil, errors.New("unable to decrypt DEK")
	}
	return plain, nil
}