package crypto

import (
	"errors"
	"time"

	"github.com/fernet/fernet-go"
)

// FernetEncryptorConfig is an interface to declate fernet key provider method.
type FernetEncryptorConfig interface {
	FernetKey() string
}

// FernetEncryptor is a structure to perform encryption/decryption functionality.
type FernetEncryptor struct {
	key       *fernet.Key
	pivotDate time.Time
}

// New is a builder function for FernetEncryptor.
func NewFernet(cfg FernetEncryptorConfig) (*FernetEncryptor, error) {
	key, err := fernet.DecodeKey(cfg.FernetKey())
	if err != nil {
		return nil, errors.New("failed to decode key")
	}

	return &FernetEncryptor{
		key:       key,
		pivotDate: time.Date(2024, 11, 7, 0, 0, 0, 0, time.UTC),
	}, nil
}

// Encrypt is an encryption function.
func (f *FernetEncryptor) Encrypt(content []byte) ([]byte, error) {
	encrypted, err := fernet.EncryptAndSign(content, f.key)
	if err != nil {
		return nil, errors.New("failed to encrypt content")
	}

	return encrypted, nil
}

// Decrypt is an decryption function.
func (f *FernetEncryptor) Decrypt(encrypted []byte) ([]byte, error) {
	decrypted := fernet.VerifyAndDecrypt(
		encrypted,
		time.Since(f.pivotDate),
		[]*fernet.Key{f.key},
	)

	if len(decrypted) == 0 {
		return nil, errors.New("failed to decrypt content")
	}

	return decrypted, nil
}
