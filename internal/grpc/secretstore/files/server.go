package files

import (
	"context"
	"errors"
	"gopasskeeper/internal/domain/models"
	"gopasskeeper/internal/grpc/interceptors"
	filesv1 "gopasskeeper/internal/grpc/secretstore/files/gen/files"
	"gopasskeeper/internal/services/secretstore/files"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	filesv1.UnimplementedFilesServer
	files Files
}

// Files is an interface for files service.
type Files interface {
	Add(
		ctx context.Context,
		uid string,
		name string,
		content []byte,
	) (*models.Message, error)
	GetSecret(
		ctx context.Context,
		uid string,
		fileID string,
	) (*models.FileSecret, error)
	Search(
		ctx context.Context,
		uid string,
		schema *models.FileSearchRequest,
	) (*models.FileSearchResponse, error)
	Remove(
		ctx context.Context,
		uid string,
		fileID string,
	) (*models.Message, error)
}

func RegisterFiles(gRPCServer *grpc.Server, files Files) {
	filesv1.RegisterFilesServer(gRPCServer, &serverAPI{files: files})
}

// Add is a method to add a file record.
func (s *serverAPI) Add(
	ctx context.Context,
	in *filesv1.FileAddRequest,
) (*filesv1.FileAddResponse, error) {
	uid, err := interceptors.ExtractUID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed extract uid")
	}

	msg, err := s.files.Add(ctx, uid, in.GetName(), in.GetContent())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed add file")
	}

	return &filesv1.FileAddResponse{Status: msg.Status, Msg: msg.Msg}, nil
}

// GetSecret is a method to get a file secret.
func (s *serverAPI) GetSecret(
	ctx context.Context,
	in *filesv1.FileSecretRequest,
) (*filesv1.FileSecretResponse, error) {
	if in.GetSecretId() == "" {
		return nil, status.Error(codes.InvalidArgument, "secret_id is required")
	}

	uid, err := interceptors.ExtractUID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed extract uid")
	}

	fileSecret, err := s.files.GetSecret(ctx, uid, in.GetSecretId())
	if err != nil {
		if errors.Is(err, files.ErrFileNotFound) {
			return nil, status.Error(codes.InvalidArgument, "failed to find file")
		}

		return nil, status.Error(codes.Internal, "failed to get secret")
	}

	return &filesv1.FileSecretResponse{
		Name:    fileSecret.Name,
		Content: fileSecret.Content,
	}, nil
}

// Search is a method to search among files records.
func (s *serverAPI) Search(
	ctx context.Context,
	in *filesv1.FileSearchRequest,
) (*filesv1.FileSearchResponse, error) {
	if in.GetLimit() == 0 {
		return nil, status.Error(codes.InvalidArgument, "limit can't be 0")
	}

	uid, err := interceptors.ExtractUID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed extract uid")
	}

	searchResponse, err := s.files.Search(
		ctx,
		uid,
		&models.FileSearchRequest{
			Substring: in.GetSubstring(),
			Offset:    uint64(in.GetOffset()),
			Limit:     uint32(in.GetLimit()),
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to search files")
	}

	items := []*filesv1.FileSearchItem{}
	for _, i := range searchResponse.Items {
		items = append(items, &filesv1.FileSearchItem{
			Id:   i.ID,
			Name: i.Name,
		})
	}

	return &filesv1.FileSearchResponse{
		Count: int64(searchResponse.Count),
		Items: items,
	}, nil
}

// Remove is a method to remove a file record.
func (s *serverAPI) Remove(
	ctx context.Context,
	in *filesv1.FileRemoveRequest,
) (*filesv1.FileRemoveResponse, error) {
	if in.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "secret_id is required")
	}

	uid, err := interceptors.ExtractUID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed extract uid")
	}

	msg, err := s.files.Remove(ctx, uid, in.GetId())
	if err != nil {
		if errors.Is(err, files.ErrFileNotFound) {
			return nil, status.Error(codes.InvalidArgument, "failed to find file")
		}
		return nil, status.Error(codes.Internal, "failed to get file")
	}

	return &filesv1.FileRemoveResponse{
		Status: msg.Status,
		Msg:    msg.Msg,
	}, nil
}
