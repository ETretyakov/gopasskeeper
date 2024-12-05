package sync

import (
	"context"
	"gopasskeeper/internal/logger"

	"github.com/pkg/errors"
)

// SyncStorage is an interface to declare sync storage methods.
type SyncStorage interface {
	Get(ctx context.Context, uid string) (string, error)
}

// Sync is a structure to describe sync service.
type Sync struct {
	log         *logger.GRPCLogger
	syncStorage SyncStorage
}

// New is a builder function for Sync.
func New(syncStorage SyncStorage) *Sync {
	return &Sync{
		log:         logger.NewGRPCLogger("sync"),
		syncStorage: syncStorage,
	}
}

// Get is a Sync method to get timestamp of recent changes.
func (s *Sync) Get(
	ctx context.Context,
	uid string,
) (string, error) {
	const op = "Sync.Get"
	log := s.log.WithOperator(op)
	log.Info("sync storage", "uid", uid)

	timestamp, err := s.syncStorage.Get(ctx, uid)
	if err != nil {
		log.Error("failed to sync storage", err)
		return "", errors.Wrap(err, op)
	}

	return timestamp, nil
}
