package notes

import (
	"context"
	"gopasskeeper/internal/domain/models"
	"gopasskeeper/internal/lib/crypto"
	"gopasskeeper/internal/logger"
	"gopasskeeper/internal/mocks"
	"gopasskeeper/internal/repository"
	"reflect"
	"testing"
)

func TestNotes_Add(t *testing.T) {
	ctx := context.Background()

	const id = "d89b92df-44e8-4d66-857a-bf7ec0a61556"
	mockedDB := mocks.NewDB(t).
		NoteAddMockedDB(id).
		AddSyncMocks()

	repo := repository.New(mockedDB.Get())
	log := logger.NewGRPCLogger("accounts-test")
	fernetEncryptor := mocks.NewFernet(t)

	type fields struct {
		log             *logger.GRPCLogger
		fernetEncryptor *crypto.FernetEncryptor
		noteStorage     NoteStorage
		syncStorage     SyncStorage
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
		want    *models.Message
		wantErr bool
	}{
		{
			fields: fields{
				log:             log,
				fernetEncryptor: fernetEncryptor,
				noteStorage:     repo.Notes,
				syncStorage:     repo.Sync,
			},
			args: args{
				ctx:     ctx,
				uid:     "31487452-31d9-4b1f-a7f8-c00b43372730",
				name:    "note",
				content: "content",
			},
			want: &models.Message{
				Status: true,
				Msg:    "Note added: note id - " + id,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Notes{
				log:             tt.fields.log,
				fernetEncryptor: tt.fields.fernetEncryptor,
				noteStorage:     tt.fields.noteStorage,
				syncStorage:     tt.fields.syncStorage,
			}
			got, err := c.Add(tt.args.ctx, tt.args.uid, tt.args.name, tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("Notes.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Notes.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotes_GetSecret(t *testing.T) {
	ctx := context.Background()

	const (
		name    = "note"
		content = "content"
	)

	fernetEncryptor := mocks.NewFernet(t)
	encContent, _ := fernetEncryptor.Encrypt([]byte(content))

	mockedDB := mocks.NewDB(t).
		NoteGetSecretMockedDB(name, string(encContent[:]))

	repo := repository.New(mockedDB.Get())
	log := logger.NewGRPCLogger("accounts-test")

	type fields struct {
		log             *logger.GRPCLogger
		fernetEncryptor *crypto.FernetEncryptor
		noteStorage     NoteStorage
		syncStorage     SyncStorage
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
			name: "Success",
			fields: fields{
				log:             log,
				fernetEncryptor: fernetEncryptor,
				noteStorage:     repo.Notes,
				syncStorage:     repo.Sync,
			},
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
			c := &Notes{
				log:             tt.fields.log,
				fernetEncryptor: tt.fields.fernetEncryptor,
				noteStorage:     tt.fields.noteStorage,
				syncStorage:     tt.fields.syncStorage,
			}
			got, err := c.GetSecret(tt.args.ctx, tt.args.uid, tt.args.noteID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Notes.GetSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Notes.GetSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotes_Search(t *testing.T) {
	ctx := context.Background()
	const (
		id   = "d89b92df-44e8-4d66-857a-bf7ec0a61556"
		name = "note"
	)

	mockedDB := mocks.NewDB(t).
		NoteSearchMockedDB(id, name)

	repo := repository.New(mockedDB.Get())
	log := logger.NewGRPCLogger("accounts-test")
	fernetEncryptor := mocks.NewFernet(t)

	type fields struct {
		log             *logger.GRPCLogger
		fernetEncryptor *crypto.FernetEncryptor
		noteStorage     NoteStorage
		syncStorage     SyncStorage
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
			fields: fields{
				log:             log,
				fernetEncryptor: fernetEncryptor,
				noteStorage:     repo.Notes,
				syncStorage:     repo.Sync,
			},
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
			c := &Notes{
				log:             tt.fields.log,
				fernetEncryptor: tt.fields.fernetEncryptor,
				noteStorage:     tt.fields.noteStorage,
				syncStorage:     tt.fields.syncStorage,
			}
			got, err := c.Search(tt.args.ctx, tt.args.uid, tt.args.schema)
			if (err != nil) != tt.wantErr {
				t.Errorf("Notes.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Notes.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotes_Remove(t *testing.T) {
	ctx := context.Background()

	mockedDB := mocks.NewDB(t).
		NoteRemoveMockedDB().
		AddSyncMocks()

	repo := repository.New(mockedDB.Get())
	log := logger.NewGRPCLogger("accounts-test")
	fernetEncryptor := mocks.NewFernet(t)

	type fields struct {
		log             *logger.GRPCLogger
		fernetEncryptor *crypto.FernetEncryptor
		noteStorage     NoteStorage
		syncStorage     SyncStorage
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
		want    *models.Message
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				log:             log,
				fernetEncryptor: fernetEncryptor,
				noteStorage:     repo.Notes,
				syncStorage:     repo.Sync,
			},
			args: args{
				ctx:    ctx,
				uid:    "31487452-31d9-4b1f-a7f8-c00b43372730",
				noteID: "d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
			want: &models.Message{
				Status: true,
				Msg:    "Note removed: note id - d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Notes{
				log:             tt.fields.log,
				fernetEncryptor: tt.fields.fernetEncryptor,
				noteStorage:     tt.fields.noteStorage,
				syncStorage:     tt.fields.syncStorage,
			}
			got, err := c.Remove(tt.args.ctx, tt.args.uid, tt.args.noteID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Notes.Remove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Notes.Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}
