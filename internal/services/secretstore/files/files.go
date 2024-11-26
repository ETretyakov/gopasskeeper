package files

import (
	"bytes"
	"context"
	"gopasskeeper/internal/config"
	"gopasskeeper/internal/domain/models"
	"gopasskeeper/internal/lib/crypto"
	"gopasskeeper/internal/logger"
	"io"

	"github.com/pkg/errors"
)

var (
	// ErrFileNotFound is an error variable to define file not found errors.
	ErrFileNotFound = errors.New("file hasn't been found")
)

// FileStorage is an interface to describe file storage methods.
type FileStorage interface {
	Add(
		ctx context.Context,
		uid string,
		name string,
		meta string,
	) (string, error)
	GetSecret(
		ctx context.Context,
		uid string,
		fileID string,
	) (string, string, error)
	Search(
		ctx context.Context,
		uid string,
		schema *models.FileSearchRequest,
	) (*models.FileSearchResponse, error)
	Remove(
		ctx context.Context,
		uid string,
		fileID string,
	) error
}

// S3Client is an interface to describe s3 client methods.
type S3Client interface {
	PutObject(ctx context.Context, name string, obj io.Reader, size int64) error
	GetObject(ctx context.Context, name string) ([]byte, error)
	RemoveObject(ctx context.Context, name string) error
}

// SyncStorage is an interface to descrube sync methods.
type SyncStorage interface {
	Set(ctx context.Context, uid string) error
}

// Files is a structure to define files service.
type Files struct {
	log             *logger.GRPCLogger
	aesEncryptor    *crypto.AESEncryptor
	fernetEncryptor *crypto.FernetEncryptor
	fileStorage     FileStorage
	s3Client        S3Client
	syncStorage     SyncStorage
}

// New is a builder function for files service.
func New(
	cfg *config.SecurityConfig,
	s3Client S3Client,
	fileStorage FileStorage,
	syncStorage SyncStorage,
) (*Files, error) {
	aesEncryptor := crypto.NewAESEncryptor(cfg)

	return &Files{
		log:          logger.NewGRPCLogger("files"),
		aesEncryptor: aesEncryptor,
		fileStorage:  fileStorage,
		s3Client:     s3Client,
		syncStorage:  syncStorage,
	}, nil
}

// Add is Files method to add file record.
func (c *Files) Add(
	ctx context.Context,
	uid string,
	name string,
	content []byte,
	meta string,
) (*models.Message, error) {
	const op = "Files.Add"
	log := c.log.WithOperator(op)
	log.Info("adding file", "uid", uid)

	encContent, err := c.aesEncryptor.Encrypt(content)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt file")
	}

	reader := bytes.NewReader(encContent)
	objName := uid + "/" + name
	if err := c.s3Client.PutObject(ctx, objName, reader, reader.Size()); err != nil {
		return nil, errors.Wrap(err, "failed save file")
	}

	encMeta, err := c.fernetEncryptor.Encrypt([]byte(meta))
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt meta")
	}

	fileID, err := c.fileStorage.Add(
		ctx,
		uid,
		name,
		string(encMeta[:]),
	)
	if err != nil {
		log.Error("failed to save file", err)
		return nil, errors.Wrap(err, op)
	}

	if err := c.syncStorage.Set(ctx, uid); err != nil {
		log.Error("failed to save file", err)
		return nil, errors.Wrap(err, op)
	}

	return &models.Message{
		Status: true,
		Msg:    "File added: file id - " + fileID,
	}, nil
}

// GetSecret is Files method to get file record with the secret.
func (c *Files) GetSecret(
	ctx context.Context,
	uid string,
	fileID string,
) (*models.FileSecret, error) {
	const op = "Files.GetSecret"
	log := c.log.WithOperator(op)
	log.Info("getting secret", "uid", uid)

	name, meta, err := c.fileStorage.GetSecret(ctx, uid, fileID)
	if err != nil {
		log.Error("failed to get file secret", err)
		return nil, errors.Wrap(err, op)
	}

	objName := uid + "/" + name
	content, err := c.s3Client.GetObject(ctx, objName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve file")
	}

	decContent, err := c.aesEncryptor.Decrypt(content)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt pin")
	}

	decMeta, err := c.fernetEncryptor.Decrypt([]byte(meta))
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt meta")
	}

	return &models.FileSecret{
		Name:    name,
		Content: decContent,
		Meta:    string(decMeta[:]),
	}, nil
}

// Search is Files method to search file records.
func (c *Files) Search(
	ctx context.Context,
	uid string,
	schema *models.FileSearchRequest,
) (*models.FileSearchResponse, error) {
	const op = "Files.Search"
	log := c.log.WithOperator(op)
	log.Info("searching items", "uid", uid)

	fileSearchResponse, err := c.fileStorage.Search(ctx, uid, schema)
	if err != nil {
		log.Error("failed to search files", err)
		return nil, errors.Wrap(err, op)
	}

	return fileSearchResponse, nil
}

// Remove is Files method to remove an file record.
func (c *Files) Remove(
	ctx context.Context,
	uid string,
	fileID string,
) (*models.Message, error) {
	const op = "Files.Search"
	log := c.log.WithOperator(op)
	log.Info("removing record", "uid", uid)

	name, _, err := c.fileStorage.GetSecret(ctx, uid, fileID)
	if err != nil {
		log.Error("failed to get file secret", err)
		return nil, errors.Wrap(err, op)
	}

	objName := uid + "/" + name
	if err := c.s3Client.RemoveObject(ctx, objName); err != nil {
		return nil, errors.Wrap(err, "failed to remove file")
	}

	if err := c.fileStorage.Remove(ctx, uid, fileID); err != nil {
		log.Error("failed to remove file", err)
		return nil, errors.Wrap(err, op)
	}

	if err := c.syncStorage.Set(ctx, uid); err != nil {
		log.Error("failed to save file", err)
		return nil, errors.Wrap(err, op)
	}

	return &models.Message{
		Status: true,
		Msg:    "File removed: file id - " + fileID,
	}, nil
}
