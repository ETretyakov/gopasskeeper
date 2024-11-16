package notes

import (
	"context"
	"gopasskeeper/internal/domain/models"
	"gopasskeeper/internal/mocks"
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestNotesRepoImpl_Add(t *testing.T) {
	ctx := context.Background()

	const id = "d89b92df-44e8-4d66-857a-bf7ec0a61556"

	mockedDB := mocks.NewDB(t).
		NoteAddMockedDB(id)

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
			fields: fields{db: mockedDB.Get()},
			args: args{
				ctx:     ctx,
				uid:     "31487452-31d9-4b1f-a7f8-c00b43372730",
				name:    "note",
				content: "content",
			},
			want:    id,
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

	const (
		name    = "note"
		content = "content"
	)
	mockedDB := mocks.NewDB(t).
		NoteGetSecretMockedDB(name, content)

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
			fields: fields{db: mockedDB.Get()},
			args: args{
				ctx:    ctx,
				uid:    "31487452-31d9-4b1f-a7f8-c00b43372730",
				noteID: "d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
			want: &models.NoteSecret{
				Name:    name,
				Content: content,
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

	const (
		id   = "d89b92df-44e8-4d66-857a-bf7ec0a61556"
		name = "note"
	)

	mockedDB := mocks.NewDB(t).
		NoteSearchMockedDB(id, name)

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
			fields: fields{db: mockedDB.Get()},
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

	mockedDB := mocks.NewDB(t).
		NoteRemoveMockedDB()

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
			fields: fields{db: mockedDB.Get()},
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
