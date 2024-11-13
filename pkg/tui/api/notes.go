package api

import (
	"context"
	"fmt"
	notesv1 "gopasskeeper/internal/grpc/secretstore/notes/gen/notes"

	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

func (a *API) SearchNote(
	substring string,
	offset uint64,
	limit uint32,
) (*notesv1.NoteSearchResponse, error) {
	md := metadata.New(map[string]string{"authorization": a.states.GetToken()})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if a.client == nil {
		return nil, nil
	}

	a.states.SetQuery(substring, offset, limit)

	resp, err := a.client.NotesAPI.Search(ctx, &notesv1.NoteSearchRequest{
		Substring: substring,
		Offset:    int64(offset),
		Limit:     int32(limit),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to search notes")
	}

	return resp, nil
}

func (a *API) GetNote(uuid string) (string, error) {
	if a.client == nil {
		return "", nil
	}

	md := metadata.New(map[string]string{"authorization": a.states.GetToken()})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := a.client.NotesAPI.GetSecret(
		ctx,
		&notesv1.NoteSecretRequest{
			SecretId: uuid,
		},
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to get note secret")
	}

	return fmt.Sprintf(
		"NAME: %s\nCONTENT: %s",
		resp.GetName(),
		resp.GetContent(),
	), nil
}

func (a *API) AddNote(name, content string) error {
	if a.client == nil {
		return nil
	}

	md := metadata.New(map[string]string{"authorization": a.states.GetToken()})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := a.client.NotesAPI.Add(
		ctx,
		&notesv1.NoteAddRequest{
			Name:    name,
			Content: content,
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to add note secret")
	}

	return nil
}

func (a *API) RemoveNote(secredID string) error {
	if a.client == nil {
		return nil
	}

	md := metadata.New(map[string]string{"authorization": a.states.GetToken()})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := a.client.NotesAPI.Remove(
		ctx,
		&notesv1.NoteRemoveRequest{
			Id: secredID,
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to remove note secret")
	}

	return nil
}
