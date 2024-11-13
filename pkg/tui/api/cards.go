package api

import (
	"context"
	"fmt"
	cardsv1 "gopasskeeper/internal/grpc/secretstore/cards/gen/cards"

	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

func (a *API) SearchCard(
	substring string,
	offset uint64,
	limit uint32,
) (*cardsv1.CardSearchResponse, error) {
	md := metadata.New(map[string]string{"authorization": a.states.GetToken()})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if a.client == nil {
		return nil, nil
	}

	a.states.SetQuery(substring, offset, limit)

	resp, err := a.client.CardsAPI.Search(ctx, &cardsv1.CardSearchRequest{
		Substring: substring,
		Offset:    int64(offset),
		Limit:     int32(limit),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to search cards")
	}

	return resp, nil
}

func (a *API) GetCard(uuid string) (string, error) {
	if a.client == nil {
		return "", nil
	}

	md := metadata.New(map[string]string{"authorization": a.states.GetToken()})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := a.client.CardsAPI.GetSecret(
		ctx,
		&cardsv1.CardSecretRequest{
			SecretId: uuid,
		},
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to get card secret")
	}

	return fmt.Sprintf(
		"NAME: %s\nNUMBER: %s\nMONTH: %d\nYEAR: %d\nCVC: %s\nPIN:%s",
		resp.GetName(),
		resp.GetNumber(),
		resp.GetMonth(),
		resp.GetYear(),
		resp.GetCvc(),
		resp.GetPin(),
	), nil
}

func (a *API) AddCard(name, number string, month, year int32, ccv, pin string) error {
	if a.client == nil {
		return nil
	}

	md := metadata.New(map[string]string{"authorization": a.states.GetToken()})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := a.client.CardsAPI.Add(
		ctx,
		&cardsv1.CardAddRequest{
			Name:   name,
			Number: number,
			Month:  month,
			Year:   year,
			Cvc:    ccv,
			Pin:    pin,
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to add card secret")
	}

	return nil
}

func (a *API) RemoveCard(secredID string) error {
	if a.client == nil {
		return nil
	}

	md := metadata.New(map[string]string{"authorization": a.states.GetToken()})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := a.client.CardsAPI.Remove(
		ctx,
		&cardsv1.CardRemoveRequest{
			Id: secredID,
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to remove card secret")
	}

	return nil
}
