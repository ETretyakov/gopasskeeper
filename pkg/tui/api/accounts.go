package api

import (
	"context"
	"fmt"
	accountsv1 "gopasskeeper/internal/grpc/secretstore/accounts/gen/accounts"

	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

func (a *API) SearchAccount(
	substring string,
	offset uint64,
	limit uint32,
) (*accountsv1.AccountSearchResponse, error) {
	md := metadata.New(map[string]string{"authorization": a.states.GetToken()})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if a.client == nil {
		return nil, nil
	}

	a.states.SetQuery(substring, offset, limit)

	resp, err := a.client.AccountsAPI.Search(ctx, &accountsv1.AccountSearchRequest{
		Substring: substring,
		Offset:    int64(offset),
		Limit:     int32(limit),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to search accounts")
	}

	return resp, nil
}

func (a *API) GetAccount(uuid string) (string, error) {
	if a.client == nil {
		return "", nil
	}

	md := metadata.New(map[string]string{"authorization": a.states.GetToken()})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := a.client.AccountsAPI.GetSecret(
		ctx,
		&accountsv1.AccountSecretRequest{
			SecretId: uuid,
		},
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to get account secret")
	}

	return fmt.Sprintf(
		"SERVICE: %s\nLOGIN: %s\nPASSWROD: %s\nMETA: %s",
		resp.GetServer(),
		resp.GetLogin(),
		resp.GetPassword(),
		resp.GetMeta(),
	), nil
}

func (a *API) AddAccount(server, login, password, meta string) error {
	if a.client == nil {
		return nil
	}

	md := metadata.New(map[string]string{"authorization": a.states.GetToken()})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := a.client.AccountsAPI.Add(
		ctx,
		&accountsv1.AccountAddRequest{
			Login:    login,
			Server:   server,
			Password: password,
			Meta:     meta,
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to add account secret")
	}

	return nil
}

func (a *API) RemoveAccount(secredID string) error {
	if a.client == nil {
		return nil
	}

	md := metadata.New(map[string]string{"authorization": a.states.GetToken()})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := a.client.AccountsAPI.Remove(
		ctx,
		&accountsv1.AccountRemoveRequest{
			Id: secredID,
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to remove account secret")
	}

	return nil
}
