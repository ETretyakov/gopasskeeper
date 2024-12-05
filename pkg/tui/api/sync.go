package api

import (
	"context"
	syncv1 "gopasskeeper/internal/grpc/sync/gen/sync"
	"os"

	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

func (a *API) GetSync() (*syncv1.SyncGetResponse, error) {
	md := metadata.New(map[string]string{"authorization": a.states.GetToken()})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if a.client == nil {
		return nil, nil
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get hostname")
	}

	resp, err := a.client.SyncAPI.Get(ctx, &syncv1.SyncGetRequest{
		DeviceId: hostname,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to search sync")
	}

	return resp, nil
}

func (a *API) Outdated() (bool, error) {
	resp, err := a.GetSync()
	if err != nil {
		return false, errors.Wrap(err, "failed to get sync")
	}

	if resp.GetTimestamp() != a.states.GetTimestamp() {
		a.states.SetTimestamp(resp.GetTimestamp())
		return true, nil
	}

	return false, nil
}
