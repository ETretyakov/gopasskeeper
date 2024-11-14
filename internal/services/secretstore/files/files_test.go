package files

import (
	"bytes"
	"context"
	"gopasskeeper/internal/config"
	"gopasskeeper/internal/domain/models"
	"gopasskeeper/internal/lib/crypto"
	"gopasskeeper/internal/logger"
	"gopasskeeper/internal/repository"
	"io"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

type MockedS3Client struct{}

func (m *MockedS3Client) PutObject(ctx context.Context, name string, obj io.Reader, size int64) error {
	return nil
}

func (m *MockedS3Client) GetObject(ctx context.Context, name string) ([]byte, error) {
	return []byte{}, nil
}

func (m *MockedS3Client) RemoveObject(ctx context.Context, name string) error {
	return nil
}

func TestFiles_Add(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rows := mock.
		NewRows([]string{"id"}).
		AddRow("d89b92df-44e8-4d66-857a-bf7ec0a61556")

	query := `
	INSERT INTO public\.sec_files\(uid, name\)
	VALUES \(.+?, .+?\)
	RETURNING id;
	`
	mock.ExpectPrepare(query)
	mock.ExpectQuery(query).WillReturnRows(rows)

	syncQuery := `
	INSERT INTO syn_timestamps\(uid, timestamp\)
	VALUES \(.+?, now\(\)\)
	ON CONFLICT ON CONSTRAINT syn_timestamps_pk
	DO UPDATE SET timestamp = excluded\.timestamp
	`
	mock.ExpectPrepare(syncQuery)
	mock.ExpectQuery(syncQuery).WillReturnRows()

	aesEncryptor := crypto.NewAESEncryptor(&config.SecurityConfig{
		AES: "3c730a7367964abd9187df2bb174d36b",
	})
	repo := repository.New(sqlxDB)
	log := logger.NewGRPCLogger("accounts-test")
	s3Client := &MockedS3Client{}

	type fields struct {
		log          *logger.GRPCLogger
		aesEncryptor *crypto.AESEncryptor
		fileStorage  FileStorage
		s3Client     S3Client
		syncStorage  SyncStorage
	}
	type args struct {
		ctx     context.Context
		uid     string
		name    string
		content []byte
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
				ctx:     ctx,
				uid:     "31487452-31d9-4b1f-a7f8-c00b43372730",
				name:    "file.txt",
				content: []byte{},
			},
			want: &models.Message{
				Status: true,
				Msg:    "File added: file id - d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
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
			got, err := c.Add(tt.args.ctx, tt.args.uid, tt.args.name, tt.args.content)
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
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rows := mock.
		NewRows([]string{"name"}).
		AddRow("file.txt")

	query := `
	SELECT sf\.name   AS \"name\"
	FROM sec_files sf
	WHERE sf\.uid = .+? AND 
	      sf\.id  = .+?
	LIMIT 1;
	`
	mock.ExpectPrepare(query)
	mock.ExpectQuery(query).WillReturnRows(rows)

	aesEncryptor := crypto.NewAESEncryptor(&config.SecurityConfig{
		AES: "3c730a7367964abd9187df2bb174d36b",
	})
	repo := repository.New(sqlxDB)
	log := logger.NewGRPCLogger("accounts-test")
	s3Client := &MockedS3Client{}

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
		want    *models.FileSecret
		wantErr bool
	}{
		{
			name: "Success",
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
			want: &models.FileSecret{
				Name:    "file.txt",
				Content: []byte{},
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
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rows := mock.
		NewRows([]string{"id", "name"}).
		AddRow("d89b92df-44e8-4d66-857a-bf7ec0a61556", "file.txt")

	query := `
	SELECT sf\.id   AS \"id\",
	       sf\.name AS \"name\"
	FROM sec_files sf
	WHERE sf\.uid = .+? AND
		  sf\.name ILIKE .+?
	ORDER BY sf\.name
	OFFSET .+?
	LIMIT  .+?;
	`
	mock.ExpectPrepare(query)
	mock.ExpectQuery(query).WillReturnRows(rows)

	countRows := mock.
		NewRows([]string{"count"}).
		AddRow(1)

	countQuery := `
	SELECT count\(\*\) AS \"count\"
	FROM sec_files sf
	WHERE sf\.uid = .+? AND
		  sf\.name ILIKE .+?;
	`
	mock.ExpectPrepare(countQuery)
	mock.ExpectQuery(countQuery).WillReturnRows(countRows)

	aesEncryptor := crypto.NewAESEncryptor(&config.SecurityConfig{
		AES: "3c730a7367964abd9187df2bb174d36b",
	})
	repo := repository.New(sqlxDB)
	log := logger.NewGRPCLogger("accounts-test")
	s3Client := &MockedS3Client{}

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
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rows := mock.
		NewRows([]string{"name"}).
		AddRow("file.txt")

	query := `
	SELECT sf\.name   AS \"name\"
	FROM sec_files sf
	WHERE sf\.uid = .+? AND 
	      sf\.id  = .+?
	LIMIT 1;
	`
	mock.ExpectPrepare(query)
	mock.ExpectQuery(query).WillReturnRows(rows)

	deleteQuery := `
	DELETE FROM sec_files
	WHERE uid = .+? AND
	      id  = .+?;
	`
	mock.ExpectPrepare(deleteQuery)
	mock.ExpectQuery(deleteQuery).WillReturnRows()

	syncQuery := `
	INSERT INTO syn_timestamps\(uid, timestamp\)
	VALUES \(.+?, now\(\)\)
	ON CONFLICT ON CONSTRAINT syn_timestamps_pk
	DO UPDATE SET timestamp = excluded\.timestamp
	`
	mock.ExpectPrepare(syncQuery)
	mock.ExpectQuery(syncQuery).WillReturnRows()

	aesEncryptor := crypto.NewAESEncryptor(&config.SecurityConfig{
		AES: "3c730a7367964abd9187df2bb174d36b",
	})
	repo := repository.New(sqlxDB)
	log := logger.NewGRPCLogger("accounts-test")
	s3Client := &MockedS3Client{}

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
