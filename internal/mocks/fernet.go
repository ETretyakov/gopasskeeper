package mocks

import (
	"gopasskeeper/internal/config"
	"gopasskeeper/internal/lib/crypto"
	"testing"
)

func NewFernet(t *testing.T) *crypto.FernetEncryptor {
	fernetEncryptor, err := crypto.NewFernet(
		&config.SecurityConfig{
			Fernet: "QijSv1fl9KAz733U_Rjxc2ribjQpJguYP2C5ezrQcwA=",
		},
	)
	if err != nil {
		t.FailNow()
	}

	return fernetEncryptor
}
