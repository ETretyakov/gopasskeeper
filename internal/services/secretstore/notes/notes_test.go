package notes

import (
	"context"
	"gopasskeeper/internal/config"
	"gopasskeeper/internal/domain/models"
	"gopasskeeper/internal/lib/crypto"
	"gopasskeeper/internal/logger"
	"gopasskeeper/internal/repository"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestNotes_Add(t *testing.T) {
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

	syncQuery := `
	INSERT INTO syn_timestamps\(uid, timestamp\)
	VALUES \(.+?, now\(\)\)
	ON CONFLICT ON CONSTRAINT syn_timestamps_pk
	DO UPDATE SET timestamp = excluded\.timestamp
	`
	mock.ExpectPrepare(syncQuery)
	mock.ExpectQuery(syncQuery).WillReturnRows()

	repo := repository.New(sqlxDB)
	log := logger.NewGRPCLogger("accounts-test")
	fernetEncryptor, _ := crypto.NewFernet(
		&config.SecurityConfig{
			Fernet: "QijSv1fl9KAz733U_Rjxc2ribjQpJguYP2C5ezrQcwA=",
		},
	)

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
				Msg:    "Note added: note id - d89b92df-44e8-4d66-857a-bf7ec0a61556",
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
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	fernetEncryptor, _ := crypto.NewFernet(
		&config.SecurityConfig{
			Fernet: "QijSv1fl9KAz733U_Rjxc2ribjQpJguYP2C5ezrQcwA=",
		},
	)

	content := "content"
	encContent, _ := fernetEncryptor.Encrypt([]byte(content))
	rows := mock.
		NewRows([]string{"name", "content"}).
		AddRow("note", string(encContent[:]))

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

	repo := repository.New(sqlxDB)
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
				Name:    "note",
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

	repo := repository.New(sqlxDB)
	log := logger.NewGRPCLogger("accounts-test")
	fernetEncryptor, _ := crypto.NewFernet(
		&config.SecurityConfig{
			Fernet: "QijSv1fl9KAz733U_Rjxc2ribjQpJguYP2C5ezrQcwA=",
		},
	)

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

	syncQuery := `
	INSERT INTO syn_timestamps\(uid, timestamp\)
	VALUES \(.+?, now\(\)\)
	ON CONFLICT ON CONSTRAINT syn_timestamps_pk
	DO UPDATE SET timestamp = excluded\.timestamp
	`
	mock.ExpectPrepare(syncQuery)
	mock.ExpectQuery(syncQuery).WillReturnRows()

	repo := repository.New(sqlxDB)
	log := logger.NewGRPCLogger("accounts-test")
	fernetEncryptor, _ := crypto.NewFernet(
		&config.SecurityConfig{
			Fernet: "QijSv1fl9KAz733U_Rjxc2ribjQpJguYP2C5ezrQcwA=",
		},
	)

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
