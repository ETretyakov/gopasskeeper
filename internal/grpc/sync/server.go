package sync

import (
	"context"
	"gopasskeeper/internal/grpc/interceptors"
	syncv1 "gopasskeeper/internal/grpc/sync/gen/sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	syncv1.UnimplementedSyncServer
	sync Sync
}

// Sync is an interface for sync service.
type Sync interface {
	Get(
		ctx context.Context,
		uid string,
	) (string, error)
}

func RegisterSync(gRPCServer *grpc.Server, sync Sync) {
	syncv1.RegisterSyncServer(gRPCServer, &serverAPI{sync: sync})
}

// Add is a method to add a note record.
func (s *serverAPI) Get(
	ctx context.Context,
	in *syncv1.SyncGetRequest,
) (*syncv1.SyncGetResponse, error) {
	uid, err := interceptors.ExtractUID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed extract uid")
	}

	timestamp, err := s.sync.Get(ctx, uid)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed add note")
	}

	return &syncv1.SyncGetResponse{Timestamp: timestamp}, nil
}
