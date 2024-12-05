package s3

import (
	"bytes"
	"context"
	"fmt"
	"gopasskeeper/internal/logger"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
)

// S3Config is an interface to describe s3 config methods.
type S3Config interface {
	Endpoint() string
	AccessKeyID() string
	SecretAccessKey() string
	BucketName() string
}

// S3ClientImpl is a structure for S3 client.
type S3ClientImpl struct {
	s3client   *minio.Client
	bucketName string
}

// New is a builder function for s3 client.
func New(ctx context.Context, cfg S3Config) (*S3ClientImpl, error) {
	client, err := minio.New(cfg.Endpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID(), cfg.SecretAccessKey(), ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start new s3 server: %w", err)
	}

	exists, err := client.BucketExists(ctx, cfg.BucketName())
	if err != nil {
		return nil, fmt.Errorf("failed to check if bucket exists: %w", err)
	}
	if !exists {
		logger.Debug("bucket hasn't been found")
		err = client.MakeBucket(ctx, cfg.BucketName(), minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create new bucket: %w", err)
		}
		logger.Debug("bucket has been created")
	}

	return &S3ClientImpl{s3client: client, bucketName: cfg.BucketName()}, nil
}

// PutObject is a method to put objects to s3.
func (s *S3ClientImpl) PutObject(
	ctx context.Context,
	name string,
	obj io.Reader,
	size int64,
) error {
	if _, err := s.s3client.PutObject(
		ctx,
		s.bucketName,
		name,
		obj,
		size,
		minio.PutObjectOptions{},
	); err != nil {
		return errors.Wrap(err, "failed to put object into bucket")
	}

	return nil
}

// GetObject is a method to get objects to s3.
func (s *S3ClientImpl) GetObject(ctx context.Context, name string) ([]byte, error) {
	obj, err := s.s3client.GetObject(ctx, s.bucketName, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get object")
	}

	defer func() {
		if err = obj.Close(); err != nil {
			logger.Error("fialed to close object", err)
		}
	}()

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(obj); err != nil {
		return nil, errors.Wrap(err, "failed to read obj from S3")
	}
	return buf.Bytes(), nil
}

// RemoveObject is a method to remove objects from s3.
func (s *S3ClientImpl) RemoveObject(ctx context.Context, name string) error {
	err := s.s3client.RemoveObject(ctx, s.bucketName, name, minio.RemoveObjectOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to get object")
	}

	return nil
}
