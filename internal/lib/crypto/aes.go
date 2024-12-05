package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"gopasskeeper/internal/logger"
	"io"

	"github.com/pkg/errors"
)

// AESEncryptorConfig is an interface to declate aes key provider method.
type AESEncryptorConfig interface {
	AESKey() string
}

// AESEncryptor is a structure to perform encryption/decryption functionality.
type AESEncryptor struct {
	key string
}

// NewAESEncryptor is a builder function for AESEncryptor.
func NewAESEncryptor(cfg AESEncryptorConfig) *AESEncryptor {
	return &AESEncryptor{key: cfg.AESKey()}
}

// Encrypt encrypts data and returns it as a byte slice.
func (a *AESEncryptor) Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(a.key))
	if err != nil {
		logger.Error("failed to get new cipher", err)
		return nil, errors.Wrap(err, "failed to build cipher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		logger.Error("failed to get gcm", err)
		return nil, errors.Wrap(err, "failed to create GCM")
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		logger.Error("failed to read full", err)
		return nil, errors.Wrap(err, "failed to create nonce")
	}

	encrypted := gcm.Seal(nonce, nonce, data, nil)
	return encrypted, nil
}

// Decrypt decrypts the encrypted data and returns it as a byte slice.
func (a *AESEncryptor) Decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(a.key))
	if err != nil {
		return nil, errors.Wrap(err, "failed to build cipher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create GCM")
	}

	if len(data) < gcm.NonceSize() {
		return nil, errors.Wrap(err, "ciphertext too short")
	}

	nonce, cText := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	decrypted, err := gcm.Open(nil, nonce, cText, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt")
	}

	return decrypted, nil
}
