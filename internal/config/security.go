package config

import "time"

// SecurityConfig is a configuration for security algorithms.
type SecurityConfig struct {
	TokenTTL    time.Duration `env:"TOKEN_TTL"  envDefault:"1h"`
	SignKey     string        `env:"SIGN_KEY"   envDefault:"87ccd01b-0b99-4f3a-8422-dd088f22cde0"`
	CertPath    string        `env:"CERT_PATH"  envDefault:"./certs/cert.pem"`
	CertKeyPath string        `env:"CERT_KEY_PATH"  envDefault:"./certs/key"`
	Fernet      string        `env:"FERNET_KEY" envDefault:"QijSv1fl9KAz733U_Rjxc2ribjQpJguYP2C5ezrQcwA="`
	AES         string        `env:"AES_KEY"    envDefault:"3c730a7367964abd9187df2bb174d36b"`
}

// FernetKey is a method to provide fernet key.
func (s *SecurityConfig) FernetKey() string {
	return s.Fernet
}

// AESKey is a method to provide aes key.
func (s *SecurityConfig) AESKey() string {
	return s.AES
}
