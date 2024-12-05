package mocks

import (
	"context"
	"io"
)

type MockedS3Client struct{}

func NewMockedS3Client() *MockedS3Client {
	return &MockedS3Client{}
}

func (m *MockedS3Client) PutObject(ctx context.Context, name string, obj io.Reader, size int64) error {
	return nil
}

func (m *MockedS3Client) GetObject(ctx context.Context, name string) ([]byte, error) {
	return []byte{}, nil
}

func (m *MockedS3Client) RemoveObject(ctx context.Context, name string) error {
	return nil
}
