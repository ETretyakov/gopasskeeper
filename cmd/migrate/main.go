package main

import (
	"errors"
	"fmt"
	"gopasskeeper/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.LoadConfig(config.WithEnvFile("server-config.yaml", "yaml"))
	if err != nil {
		log.Error().Msg("failed to load config")
		return
	}

	cfg.Show()

	if cfg.PostgreSQL.MigrationFolder == "" {
		panic("migrations-path is required")
	}

	m, err := migrate.New(
		"file://"+cfg.PostgreSQL.MigrationFolder,
		cfg.PostgreSQL.DSN()+"?sslmode=disable",
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		panic(err)
	}
}
