package config

import (
	"fmt"
	"time"
)

// PostgreSQLConfig is a configuration for PostgreSQL database.
type PostgreSQLConfig struct {
	Host            string        `env:"HOST"             envDefault:"127.0.0.1"`
	Port            int           `env:"PORT"             envDefault:"5432"`
	User            string        `env:"USER"             envDefault:"postgres"`
	Password        string        `env:"PASSWORD"         envDefault:"thepass123"`
	Database        string        `env:"DATABASE"         envDefault:"postgres"`
	MaxOpenConn     int           `env:"MAX_OPEN_CONN"    envDefault:"10"`
	IdleConn        int           `env:"MAX_IDLE_CONN"    envDefault:"10"`
	PingInterval    time.Duration `env:"DURATION"         envDefault:"5s"`
	MigrationFolder string        `env:"MIGRATION_FOLDER" envDefault:"./migrations"`
}

// DSN is a method which generates PostgreSQL DSN string.
func (pg *PostgreSQLConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		pg.User,
		pg.Password,
		pg.Host,
		pg.Port,
		pg.Database,
	)
}
