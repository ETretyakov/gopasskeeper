package files

import (
	"context"
	"gopasskeeper/internal/domain/models"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestFilesRepoImpl_Add(t *testing.T) {
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

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx  context.Context
		uid  string
		name string
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
			fields: fields{db: sqlxDB},
			args: args{
				ctx:  ctx,
				uid:  "31487452-31d9-4b1f-a7f8-c00b43372730",
				name: "file.txt",
			},
			want:    "d89b92df-44e8-4d66-857a-bf7ec0a61556",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FilesRepoImpl{
				db: tt.fields.db,
			}
			got, err := r.Add(tt.args.ctx, tt.args.uid, tt.args.name)
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
		want    string
		wantErr bool
	}{
		{
			name:   "Success",
			fields: fields{db: sqlxDB},
			args: args{
				ctx:    ctx,
				uid:    "31487452-31d9-4b1f-a7f8-c00b43372730",
				fileID: "d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
			want:    "file.txt",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &FilesRepoImpl{
				db: tt.fields.db,
			}
			got, err := r.GetSecret(tt.args.ctx, tt.args.uid, tt.args.fileID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilesRepoImpl.GetSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FilesRepoImpl.GetSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilesRepoImpl_Search(t *testing.T) {
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
			fields: fields{db: sqlxDB},
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
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	query := `
	DELETE FROM sec_files
	WHERE uid = .+? AND
	      id  = .+?;
	`
	mock.ExpectPrepare(query)
	mock.ExpectQuery(query).WillReturnRows()

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
			fields: fields{db: sqlxDB},
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
