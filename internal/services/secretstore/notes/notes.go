package notes

import (
	"context"
	"gopasskeeper/internal/config"
	"gopasskeeper/internal/domain/models"
	"gopasskeeper/internal/lib/crypto"
	"gopasskeeper/internal/logger"

	"github.com/pkg/errors"
)

var (
	// ErrNoteNotFound is an error variable to define note not found errors.
	ErrNoteNotFound = errors.New("note hasn't been found")
)

// NoteStorage is an interface to describe note storage methods.
type NoteStorage interface {
	Add(
		ctx context.Context,
		uid string,
		name string,
		content string,
	) (string, error)
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
	) error
}

// SyncStorage is an interface to descrube sync methods.
type SyncStorage interface {
	Set(ctx context.Context, uid string) error
}

// Notes is a structure to define notes service.
type Notes struct {
	log             *logger.GRPCLogger
	fernetEncryptor *crypto.FernetEncryptor
	noteStorage     NoteStorage
	syncStorage     SyncStorage
}

// New is a builder function for notes service.
func New(
	cfg *config.SecurityConfig,
	noteStorage NoteStorage,
	syncStorage SyncStorage,
) (*Notes, error) {
	fernetEncryptor, err := crypto.NewFernet(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build fernet encryptor")
	}

	return &Notes{
		log:             logger.NewGRPCLogger("notes"),
		fernetEncryptor: fernetEncryptor,
		noteStorage:     noteStorage,
		syncStorage:     syncStorage,
	}, nil
}

// Add is Notes method to add note record.
func (c *Notes) Add(
	ctx context.Context,
	uid string,
	name string,
	content string,
) (*models.Message, error) {
	const op = "Notes.Add"
	log := c.log.WithOperator(op)
	log.Info("adding note", "uid", uid)

	encContent, err := c.fernetEncryptor.Encrypt([]byte(content))
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt pin")
	}

	noteID, err := c.noteStorage.Add(
		ctx,
		uid,
		name,
		string(encContent[:]),
	)
	if err != nil {
		log.Error("failed to save note", err)
		return nil, errors.Wrap(err, op)
	}

	if err := c.syncStorage.Set(ctx, uid); err != nil {
		log.Error("failed to save note", err)
		return nil, errors.Wrap(err, op)
	}

	return &models.Message{
		Status: true,
		Msg:    "Note added: note id - " + noteID,
	}, nil
}

// GetSecret is Notes method to get note record with the secret.
func (c *Notes) GetSecret(
	ctx context.Context,
	uid string,
	noteID string,
) (*models.NoteSecret, error) {
	const op = "Notes.GetSecret"
	log := c.log.WithOperator(op)
	log.Info("getting secret", "uid", uid)

	noteSecret, err := c.noteStorage.GetSecret(ctx, uid, noteID)
	if err != nil {
		log.Error("failed to get note secret", err)
		return nil, errors.Wrap(err, op)
	}

	decContent, err := c.fernetEncryptor.Decrypt([]byte(noteSecret.Content))
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt pin")
	}
	noteSecret.Content = string(decContent[:])

	return noteSecret, nil
}

// Search is Notes method to search note records.
func (c *Notes) Search(
	ctx context.Context,
	uid string,
	schema *models.NoteSearchRequest,
) (*models.NoteSearchResponse, error) {
	const op = "Notes.Search"
	log := c.log.WithOperator(op)
	log.Info("searching items", "uid", uid)

	noteSearchResponse, err := c.noteStorage.Search(ctx, uid, schema)
	if err != nil {
		log.Error("failed to search notes", err)
		return nil, errors.Wrap(err, op)
	}

	return noteSearchResponse, nil
}

// Remove is Notes method to remove an note record.
func (c *Notes) Remove(
	ctx context.Context,
	uid string,
	noteID string,
) (*models.Message, error) {
	const op = "Notes.Search"
	log := c.log.WithOperator(op)
	log.Info("removing record", "uid", uid)

	if err := c.noteStorage.Remove(ctx, uid, noteID); err != nil {
		log.Error("failed to remove note", err)
		return nil, errors.Wrap(err, op)
	}

	if err := c.syncStorage.Set(ctx, uid); err != nil {
		log.Error("failed to save note", err)
		return nil, errors.Wrap(err, op)
	}

	return &models.Message{
		Status: true,
		Msg:    "Note removed: note id - " + noteID,
	}, nil
}
