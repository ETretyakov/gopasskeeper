package sync

import (
	"context"
	"gopasskeeper/internal/logger"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// SyncRepoImpl is a structure to implement sync repository.
type SyncRepoImpl struct {
	db *sqlx.DB
}

// New is a builder function for SyncRepoImpl.
func New(db *sqlx.DB) *SyncRepoImpl {
	return &SyncRepoImpl{db: db}
}

// Get is a method to get sync timestamp.
func (s *SyncRepoImpl) Get(ctx context.Context, uid string) (string, error) {
	const op = "repository.Sync.Get"

	stmt := `
	SELECT st.timestamp AS "timestamp"
	FROM syn_timestamps st
	WHERE st.uid = :uid
	LIMIT 1;`

	namedStmt, err := s.db.PrepareNamedContext(ctx, stmt)
	if err != nil {
		logger.Error("failed to prepare named context", err)
		return "", errors.Wrap(err, op)
	}

	arg := map[string]interface{}{
		"uid": uid,
	}

	var timestamp string
	if err := namedStmt.QueryRowxContext(ctx, arg).Scan(&timestamp); err != nil {
		logger.Error("failed to query row context", err)
		return "", errors.Wrap(err, op)
	}

	return timestamp, nil
}

// Set is a method to set sync timestamp.
func (s *SyncRepoImpl) Set(ctx context.Context, uid string) error {
	const op = "repository.Sync.Set"

	stmt := `
	INSERT INTO syn_timestamps(uid, timestamp)
	VALUES (:uid, now())
	ON CONFLICT ON CONSTRAINT syn_timestamps_pk
	DO UPDATE SET timestamp = excluded.timestamp
	`

	namedStmt, err := s.db.PrepareNamedContext(ctx, stmt)
	if err != nil {
		logger.Error("failed to prepare named context", err)
		return errors.Wrap(err, op)
	}

	arg := map[string]interface{}{
		"uid": uid,
	}

	if err := namedStmt.QueryRowxContext(ctx, arg).Err(); err != nil {
		logger.Error("failed to query row context", err)
		return errors.Wrap(err, op)
	}

	return nil
}
