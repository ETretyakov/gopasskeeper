package files

import (
	"context"
	"gopasskeeper/internal/domain/models"
	"gopasskeeper/internal/mocks"
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestFilesRepoImpl_Add(t *testing.T) {
	ctx := context.Background()

	const id = "d89b92df-44e8-4d66-857a-bf7ec0a61556"

	mockedDB := mocks.NewDB(t).
		FileAddMockedDB(id)

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx  context.Context
		uid  string
		name string
		meta string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "Success",
			fields: fields{db: mockedDB.Get()},
			args: args{
				ctx:  ctx,
				uid:  "31487452-31d9-4b1f-a7f8-c00b43372730",
				name: "file.txt",
				meta: "meta",
			},
			want:    id,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FilesRepoImpl{
				db: tt.fields.db,
			}
			got, err := r.Add(tt.args.ctx, tt.args.uid, tt.args.name, tt.args.meta)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilesRepoImpl.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FilesRepoImpl.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilesRepoImpl_GetSecret(t *testing.T) {
	ctx := context.Background()

	const (
		name = "file.txt"
		meta = "meta"
	)
	mockedDB := mocks.NewDB(t).
		FileGetSecretMockedDB(name, meta)

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		uid    string
		fileID string
		meta   string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantName string
		wantMeta string
		wantErr  bool
	}{
		{
			name:   "Success",
			fields: fields{db: mockedDB.Get()},
			args: args{
				ctx:    ctx,
				uid:    "31487452-31d9-4b1f-a7f8-c00b43372730",
				fileID: "d89b92df-44e8-4d66-857a-bf7ec0a61556",
				meta:   "meta",
			},
			wantName: name,
			wantMeta: meta,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FilesRepoImpl{
				db: tt.fields.db,
			}
			gotName, gotMeta, err := r.GetSecret(tt.args.ctx, tt.args.uid, tt.args.fileID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilesRepoImpl.GetSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotName != tt.wantName {
				t.Errorf("FilesRepoImpl.GetSecret() = %v, want %v", gotName, tt.wantName)
			}
			if gotMeta != tt.wantMeta {
				t.Errorf("FilesRepoImpl.GetSecret() = %v, want %v", gotMeta, tt.wantMeta)
			}
		})
	}
}

func TestFilesRepoImpl_Search(t *testing.T) {
	ctx := context.Background()

	const (
		id   = "d89b92df-44e8-4d66-857a-bf7ec0a61556"
		name = "file.txt"
	)

	mockedDB := mocks.NewDB(t).
		FileSearchMockedDB(id, name)

	type fields struct {
		db *sqlx.DB
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
			name:   "Success",
			fields: fields{db: mockedDB.Get()},
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
						ID:   id,
						Name: name,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FilesRepoImpl{
				db: tt.fields.db,
			}
			got, err := r.Search(tt.args.ctx, tt.args.uid, tt.args.schema)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilesRepoImpl.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilesRepoImpl.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilesRepoImpl_Remove(t *testing.T) {
	ctx := context.Background()

	mockedDB := mocks.NewDB(t).
		FileRemoveMockedDB()

	type fields struct {
		db *sqlx.DB
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
		wantErr bool
	}{
		{
			name:   "Success",
			fields: fields{db: mockedDB.Get()},
			args: args{
				ctx:    ctx,
				uid:    "31487452-31d9-4b1f-a7f8-c00b43372730",
				fileID: "d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FilesRepoImpl{
				db: tt.fields.db,
			}
			if err := r.Remove(tt.args.ctx, tt.args.uid, tt.args.fileID); (err != nil) != tt.wantErr {
				t.Errorf("FilesRepoImpl.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
