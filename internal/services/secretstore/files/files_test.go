package files

import (
	"bytes"
	"context"
	"gopasskeeper/internal/domain/models"
	"gopasskeeper/internal/lib/crypto"
	"gopasskeeper/internal/logger"
	"gopasskeeper/internal/mocks"
	"gopasskeeper/internal/repository"
	"reflect"
	"testing"
)

func TestFiles_Add(t *testing.T) {
	ctx := context.Background()

	const id = "d89b92df-44e8-4d66-857a-bf7ec0a61556"
	mockedDB := mocks.NewDB(t).
		FileAddMockedDB(id).
		AddSyncMocks()

	aesEncryptor := mocks.NewAESEncryptor()
	repo := repository.New(mockedDB.Get())
	log := logger.NewGRPCLogger("accounts-test")
	fernetEncryptor := mocks.NewFernet(t)
	s3Client := mocks.NewMockedS3Client()

	type fields struct {
		log             *logger.GRPCLogger
		aesEncryptor    *crypto.AESEncryptor
		fernetEncryptor *crypto.FernetEncryptor
		fileStorage     FileStorage
		s3Client        S3Client
		syncStorage     SyncStorage
	}
	type args struct {
		ctx     context.Context
		uid     string
		name    string
		content []byte
		meta    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Message
		wantErr bool
	}{
		{
			fields: fields{
				log:             log,
				aesEncryptor:    aesEncryptor,
				fernetEncryptor: fernetEncryptor,
				fileStorage:     repo.Files,
				s3Client:        s3Client,
				syncStorage:     repo.Sync,
			},
			args: args{
				ctx:     ctx,
				uid:     "31487452-31d9-4b1f-a7f8-c00b43372730",
				name:    "file.txt",
				content: []byte{},
				meta:    "meta",
			},
			want: &models.Message{
				Status: true,
				Msg:    "File added: file id - " + id,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Files{
				log:             tt.fields.log,
				aesEncryptor:    tt.fields.aesEncryptor,
				fernetEncryptor: tt.fields.fernetEncryptor,
				fileStorage:     tt.fields.fileStorage,
				s3Client:        tt.fields.s3Client,
				syncStorage:     tt.fields.syncStorage,
			}
			got, err := c.Add(tt.args.ctx, tt.args.uid, tt.args.name, tt.args.content, tt.args.meta)
			if (err != nil) != tt.wantErr {
				t.Errorf("Files.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Files.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFiles_GetSecret(t *testing.T) {
	ctx := context.Background()

	const (
		name = "file.txt"
		meta = "meta"
	)

	log := logger.NewGRPCLogger("accounts-test")
	s3Client := mocks.NewMockedS3Client()
	fernetEncryptor := mocks.NewFernet(t)
	aesEncryptor := mocks.NewAESEncryptor()

	encMeta, _ := fernetEncryptor.Encrypt([]byte(meta))

	mockedDB := mocks.NewDB(t).
		FileGetSecretMockedDB(name, string(encMeta[:]))

	repo := repository.New(mockedDB.Get())

	type fields struct {
		log             *logger.GRPCLogger
		aesEncryptor    *crypto.AESEncryptor
		fernetEncryptor *crypto.FernetEncryptor
		fileStorage     FileStorage
		s3Client        S3Client
		syncStorage     SyncStorage
	}
	type args struct {
		ctx    context.Context
		uid    string
		fileID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.FileSecret
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				log:             log,
				aesEncryptor:    aesEncryptor,
				fernetEncryptor: fernetEncryptor,
				fileStorage:     repo.Files,
				s3Client:        s3Client,
				syncStorage:     repo.Sync,
			},
			args: args{
				ctx:    ctx,
				uid:    "31487452-31d9-4b1f-a7f8-c00b43372730",
				fileID: "d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
			want: &models.FileSecret{
				Name:    "file.txt",
				Content: []byte{},
				Meta:    "meta",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Files{
				log:             tt.fields.log,
				aesEncryptor:    tt.fields.aesEncryptor,
				fernetEncryptor: tt.fields.fernetEncryptor,
				fileStorage:     tt.fields.fileStorage,
				s3Client:        tt.fields.s3Client,
				syncStorage:     tt.fields.syncStorage,
			}
			got, err := c.GetSecret(tt.args.ctx, tt.args.uid, tt.args.fileID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Files.GetSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !bytes.Equal(got.Content, tt.want.Content) {
				t.Errorf("Files.GetSecret() = %v, want %v", got, tt.want)
			}
			if got.Name != tt.want.Name {
				t.Errorf("Files.GetSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFiles_Search(t *testing.T) {
	ctx := context.Background()
	const (
		id   = "d89b92df-44e8-4d66-857a-bf7ec0a61556"
		name = "file.txt"
	)

	mockedDB := mocks.NewDB(t).
		FileSearchMockedDB(id, name)

	repo := repository.New(mockedDB.Get())

	aesEncryptor := mocks.NewAESEncryptor()
	log := logger.NewGRPCLogger("accounts-test")
	s3Client := mocks.NewMockedS3Client()

	type fields struct {
		log          *logger.GRPCLogger
		aesEncryptor *crypto.AESEncryptor
		fileStorage  FileStorage
		s3Client     S3Client
		syncStorage  SyncStorage
	}
	type args struct {
		ctx    context.Context
		uid    string
		schema *models.FileSearchRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.FileSearchResponse
		wantErr bool
	}{
		{
			fields: fields{
				log:          log,
				aesEncryptor: aesEncryptor,
				fileStorage:  repo.Files,
				s3Client:     s3Client,
				syncStorage:  repo.Sync,
			},
			args: args{
				ctx: ctx,
				uid: "31487452-31d9-4b1f-a7f8-c00b43372730",
				schema: &models.FileSearchRequest{
					Substring: "",
					Offset:    0,
					Limit:     100,
				},
			},
			want: &models.FileSearchResponse{
				Count: 1,
				Items: []*models.FileSearchItem{
					{
						ID:   "d89b92df-44e8-4d66-857a-bf7ec0a61556",
						Name: "file.txt",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Files{
				log:          tt.fields.log,
				aesEncryptor: tt.fields.aesEncryptor,
				fileStorage:  tt.fields.fileStorage,
				s3Client:     tt.fields.s3Client,
				syncStorage:  tt.fields.syncStorage,
			}
			got, err := c.Search(tt.args.ctx, tt.args.uid, tt.args.schema)
			if (err != nil) != tt.wantErr {
				t.Errorf("Files.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Files.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFiles_Remove(t *testing.T) {
	ctx := context.Background()

	const (
		name = "file.txt"
		meta = "meta"
	)

	mockedDB := mocks.NewDB(t).
		FileGetSecretMockedDB(name, meta).
		FileRemoveMockedDB().
		AddSyncMocks()

	repo := repository.New(mockedDB.Get())

	aesEncryptor := mocks.NewAESEncryptor()
	log := logger.NewGRPCLogger("accounts-test")
	s3Client := mocks.NewMockedS3Client()

	type fields struct {
		log          *logger.GRPCLogger
		aesEncryptor *crypto.AESEncryptor
		fileStorage  FileStorage
		s3Client     S3Client
		syncStorage  SyncStorage
	}
	type args struct {
		ctx    context.Context
		uid    string
		fileID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Message
		wantErr bool
	}{
		{
			fields: fields{
				log:          log,
				aesEncryptor: aesEncryptor,
				fileStorage:  repo.Files,
				s3Client:     s3Client,
				syncStorage:  repo.Sync,
			},
			args: args{
				ctx:    ctx,
				uid:    "31487452-31d9-4b1f-a7f8-c00b43372730",
				fileID: "d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
			want: &models.Message{
				Status: true,
				Msg:    "File removed: file id - d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Files{
				log:          tt.fields.log,
				aesEncryptor: tt.fields.aesEncryptor,
				fileStorage:  tt.fields.fileStorage,
				s3Client:     tt.fields.s3Client,
				syncStorage:  tt.fields.syncStorage,
			}
			got, err := c.Remove(tt.args.ctx, tt.args.uid, tt.args.fileID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Files.Remove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Files.Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}
