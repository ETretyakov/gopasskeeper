package postgres

import (
	"context"
	"gopasskeeper/internal/config"
	"gopasskeeper/internal/logger"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Storage is a structure to declare methods to work with database.
type Storage struct {
	DB *sqlx.DB
}

// Close is a Storage method to close database connection.
func (s *Storage) Close() error {
	if err := s.DB.Close(); err != nil {
		return errors.Wrap(err, "failed to close database connection")
	}
	return nil
}

// New is a builder function for Storage.
func New(
	ctx context.Context,
	cfg *config.PostgreSQLConfig,
) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sqlx.Open("pgx", cfg.DSN())
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	db.SetConnMaxLifetime(0)
	db.SetConnMaxIdleTime(0)
	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetMaxIdleConns(cfg.IdleConn)

	go func() {
		ticker := time.NewTicker(cfg.PingInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := db.Ping(); err != nil {
					logger.Error("error ping db", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return &Storage{DB: db}, nil
}
