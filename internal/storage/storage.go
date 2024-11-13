package storage

import (
	"context"
	"gopasskeeper/internal/closer"
	"gopasskeeper/internal/config"
	"gopasskeeper/internal/storage/postgres"

	"github.com/pkg/errors"

	_ "github.com/jackc/pgx/stdlib"
)

var (
	// ErrUserExists is an error for user exists cases.
	ErrUserExists = errors.New("user already exists")
	// ErrUserNotFound is an error for user not found cases.
	ErrUserNotFound = errors.New("user not found")
)

// NewPostgresDB is a builder function for postgres.Storage.
func NewPostgresDB(
	ctx context.Context,
	cfg *config.PostgreSQLConfig,
) (*postgres.Storage, error) {
	storage, err := postgres.New(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build postgres")
	}

	closer.Add(storage.Close)

	return storage, nil
}
