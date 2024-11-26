package api

import (
	"context"
	"fmt"
	filesv1 "gopasskeeper/internal/grpc/secretstore/files/gen/files"
	"os"

	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

func (a *API) SearchFile(
	substring string,
	offset uint64,
	limit uint32,
) (*filesv1.FileSearchResponse, error) {
	md := metadata.New(map[string]string{"authorization": a.states.GetToken()})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if a.client == nil {
		return nil, nil
	}

	a.states.SetQuery(substring, offset, limit)

	resp, err := a.client.FilesAPI.Search(ctx, &filesv1.FileSearchRequest{
		Substring: substring,
		Offset:    int64(offset),
		Limit:     int32(limit),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to search files")
	}

	return resp, nil
}

func (a *API) GetFile(uuid string, filePath string) (string, error) {
	if a.client == nil {
		return "", nil
	}

	md := metadata.New(map[string]string{"authorization": a.states.GetToken()})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return "", errors.Wrap(err, "failed to open file for writing")
	}
	defer file.Close()

	resp, err := a.client.FilesAPI.GetSecret(
		ctx,
		&filesv1.FileSecretRequest{
			SecretId: uuid,
		},
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to get file secret")
	}

	if _, err := file.Write(resp.GetContent()); err != nil {
		return "", errors.Wrap(err, "failed to write file")
	}

	return fmt.Sprintf(
		"File has been written: %s",
		filePath,
	), nil
}

func (a *API) AddFile(name, filePath, meta string) error {
	if a.client == nil {
		return nil
	}

	md := metadata.New(map[string]string{"authorization": a.states.GetToken()})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}

	_, err = a.client.FilesAPI.Add(
		ctx,
		&filesv1.FileAddRequest{
			Name:    name,
			Content: content,
			Meta:    meta,
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to add file secret")
	}

	return nil
}

func (a *API) RemoveFile(secredID string) error {
	if a.client == nil {
		return nil
	}

	md := metadata.New(map[string]string{"authorization": a.states.GetToken()})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := a.client.FilesAPI.Remove(
		ctx,
		&filesv1.FileRemoveRequest{
			Id: secredID,
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to remove file secret")
	}

	return nil
}
