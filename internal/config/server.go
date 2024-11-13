package config

import (
	"fmt"
	"time"
)

// ServerConfig is a configuration for GRPC server.
type ServerConfig struct {
	Host    string        `env:"ADDRESS" envDefault:"127.0.0.1"`
	Port    uint32        `env:"PORT"    envDefault:"9090"`
	Timeout time.Duration `env:"timeout" envDefault:"30s"`
}

// Address is a method to generate address using host and port.
func (s *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
