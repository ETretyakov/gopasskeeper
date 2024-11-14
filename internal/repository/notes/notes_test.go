package notes

import (
	"context"
	"gopasskeeper/internal/domain/models"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestNotesRepoImpl_Add(t *testing.T) {
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
	INSERT INTO public\.sec_notes\(uid, name, content\)
	VALUES \(.+?, .+?, .+?\)
	RETURNING id;
	`
	mock.ExpectPrepare(query)
	mock.ExpectQuery(query).WillReturnRows(rows)

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx     context.Context
		uid     string
		name    string
		content string
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
				ctx:     ctx,
				uid:     "31487452-31d9-4b1f-a7f8-c00b43372730",
				name:    "note",
				content: "content",
			},
			want:    "d89b92df-44e8-4d66-857a-bf7ec0a61556",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &NotesRepoImpl{
				db: tt.fields.db,
			}
			got, err := r.Add(tt.args.ctx, tt.args.uid, tt.args.name, tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("NotesRepoImpl.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NotesRepoImpl.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotesRepoImpl_GetSecret(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rows := mock.
		NewRows([]string{"name", "content"}).
		AddRow("note", "content")

	query := `
	SELECT sn\.name   AS \"name\",
		   sn\.content AS \"content\"
	FROM sec_notes sn
	WHERE sn\.uid = .+? AND 
	      sn\.id  = .+?
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
		noteID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.NoteSecret
		wantErr bool
	}{
		{
			fields: fields{db: sqlxDB},
			args: args{
				ctx:    ctx,
				uid:    "31487452-31d9-4b1f-a7f8-c00b43372730",
				noteID: "d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
			want: &models.NoteSecret{
				Name:    "note",
				Content: "content",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &NotesRepoImpl{
				db: tt.fields.db,
			}
			got, err := r.GetSecret(tt.args.ctx, tt.args.uid, tt.args.noteID)
			if (err != nil) != tt.wantErr {
				t.Errorf("NotesRepoImpl.GetSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NotesRepoImpl.GetSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotesRepoImpl_Search(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rows := mock.
		NewRows([]string{"id", "name"}).
		AddRow("d89b92df-44e8-4d66-857a-bf7ec0a61556", "note")

	query := `
	SELECT sn\.id   AS \"id\",
	       sn\.name AS \"name\"
	FROM sec_notes sn
	WHERE sn\.uid = .+? AND
		  sn\.name ILIKE .+?
	ORDER BY sn\.name
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
	FROM sec_notes sn
	WHERE sn\.uid = .+? AND
		  sn\.name ILIKE .+?;
	`
	mock.ExpectPrepare(countQuery)
	mock.ExpectQuery(countQuery).WillReturnRows(countRows)

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		uid    string
		schema *models.NoteSearchRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.NoteSearchResponse
		wantErr bool
	}{
		{
			name:   "Success",
			fields: fields{db: sqlxDB},
			args: args{
				ctx: ctx,
				uid: "31487452-31d9-4b1f-a7f8-c00b43372730",
				schema: &models.NoteSearchRequest{
					Substring: "",
					Offset:    0,
					Limit:     100,
				},
			},
			want: &models.NoteSearchResponse{
				Count: 1,
				Items: []*models.NoteSearchItem{
					{
						ID:   "d89b92df-44e8-4d66-857a-bf7ec0a61556",
						Name: "note",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &NotesRepoImpl{
				db: tt.fields.db,
			}
			got, err := r.Search(tt.args.ctx, tt.args.uid, tt.args.schema)
			if (err != nil) != tt.wantErr {
				t.Errorf("NotesRepoImpl.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NotesRepoImpl.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotesRepoImpl_Remove(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	query := `
	DELETE FROM sec_notes
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
		noteID string
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
				noteID: "d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &NotesRepoImpl{
				db: tt.fields.db,
			}
			if err := r.Remove(tt.args.ctx, tt.args.uid, tt.args.noteID); (err != nil) != tt.wantErr {
				t.Errorf("NotesRepoImpl.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
