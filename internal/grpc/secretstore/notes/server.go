package notes

import (
	"context"
	"errors"
	"gopasskeeper/internal/domain/models"
	"gopasskeeper/internal/grpc/interceptors"
	notesv1 "gopasskeeper/internal/grpc/secretstore/notes/gen/notes"
	"gopasskeeper/internal/services/secretstore/notes"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	notesv1.UnimplementedNotesServer
	notes Notes
}

// Notes is an interface for notes service.
type Notes interface {
	Add(
		ctx context.Context,
		uid string,
		name string,
		content string,
		meta string,
	) (*models.Message, error)
	GetSecret(
		ctx context.Context,
		uid string,
		noteID string,
	) (*models.NoteSecret, error)
	Search(
		ctx context.Context,
		uid string,
		schema *models.NoteSearchRequest,
	) (*models.NoteSearchResponse, error)
	Remove(
		ctx context.Context,
		uid string,
		noteID string,
	) (*models.Message, error)
}

// RegisterNotes is a function to register Notes server.
func RegisterNotes(gRPCServer *grpc.Server, notes Notes) {
	notesv1.RegisterNotesServer(gRPCServer, &serverAPI{notes: notes})
}

// Add is a method to add a note record.
func (s *serverAPI) Add(
	ctx context.Context,
	in *notesv1.NoteAddRequest,
) (*notesv1.NoteAddResponse, error) {
	uid, err := interceptors.ExtractUID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed extract uid")
	}

	if in.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "note must have name")
	}

	if in.GetContent() == "" {
		return nil, status.Error(codes.InvalidArgument, "note must have content")
	}

	msg, err := s.notes.Add(ctx, uid, in.GetName(), in.GetContent(), in.GetMeta())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed add note")
	}

	return &notesv1.NoteAddResponse{Status: msg.Status, Msg: msg.Msg}, nil
}

// GetSecret is a method to get a note secret.
func (s *serverAPI) GetSecret(
	ctx context.Context,
	in *notesv1.NoteSecretRequest,
) (*notesv1.NoteSecretResponse, error) {
	if in.GetSecretId() == "" {
		return nil, status.Error(codes.InvalidArgument, "secret_id is required")
	}

	uid, err := interceptors.ExtractUID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed extract uid")
	}

	noteSecret, err := s.notes.GetSecret(ctx, uid, in.GetSecretId())
	if err != nil {
		if errors.Is(err, notes.ErrNoteNotFound) {
			return nil, status.Error(codes.InvalidArgument, "failed to find note")
		}

		return nil, status.Error(codes.Internal, "failed to get secret")
	}

	return &notesv1.NoteSecretResponse{
		Name:    noteSecret.Name,
		Content: noteSecret.Content,
		Meta:    noteSecret.Meta,
	}, nil
}

// Search is a method to search among notes records.
func (s *serverAPI) Search(
	ctx context.Context,
	in *notesv1.NoteSearchRequest,
) (*notesv1.NoteSearchResponse, error) {
	if in.GetLimit() == 0 {
		return nil, status.Error(codes.InvalidArgument, "limit can't be 0")
	}

	uid, err := interceptors.ExtractUID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed extract uid")
	}

	searchResponse, err := s.notes.Search(
		ctx,
		uid,
		&models.NoteSearchRequest{
			Substring: in.GetSubstring(),
			Offset:    uint64(in.GetOffset()),
			Limit:     uint32(in.GetLimit()),
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to search notes")
	}

	items := []*notesv1.NoteSearchItem{}
	for _, i := range searchResponse.Items {
		items = append(items, &notesv1.NoteSearchItem{
			Id:   i.ID,
			Name: i.Name,
		})
	}

	return &notesv1.NoteSearchResponse{
		Count: int64(searchResponse.Count),
		Items: items,
	}, nil
}

// Remove is a method to remove a note record.
func (s *serverAPI) Remove(
	ctx context.Context,
	in *notesv1.NoteRemoveRequest,
) (*notesv1.NoteRemoveResponse, error) {
	if in.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "secret_id is required")
	}

	uid, err := interceptors.ExtractUID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed extract uid")
	}

	msg, err := s.notes.Remove(ctx, uid, in.GetId())
	if err != nil {
		if errors.Is(err, notes.ErrNoteNotFound) {
			return nil, status.Error(codes.InvalidArgument, "failed to find note")
		}
		return nil, status.Error(codes.Internal, "failed to get note")
	}

	return &notesv1.NoteRemoveResponse{
		Status: msg.Status,
		Msg:    msg.Msg,
	}, nil
}
