package mocks

import (
	"gopasskeeper/internal/config"
	"gopasskeeper/internal/lib/crypto"
)

func NewAESEncryptor() *crypto.AESEncryptor {
	return crypto.NewAESEncryptor(&config.SecurityConfig{
		AES: "3c730a7367964abd9187df2bb174d36b",
	})
}
