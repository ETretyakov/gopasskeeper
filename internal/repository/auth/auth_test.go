package auth

import (
	"context"
	"gopasskeeper/internal/domain/models"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestAuthRepoImpl_SaveUser(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rows := mock.
		NewRows([]string{"uid"}).
		AddRow("31487452-31d9-4b1f-a7f8-c00b43372730")

	query := `
	INSERT INTO public.usr_users\(login, pass_hash\)
	VALUES \(.+?, .+?\)
	RETURNING id;
	`
	mock.ExpectPrepare(query)
	mock.ExpectQuery(query).WillReturnRows(rows)

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx      context.Context
		login    string
		passHash []byte
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
				ctx:      ctx,
				login:    "user",
				passHash: []byte{},
			},
			want:    "31487452-31d9-4b1f-a7f8-c00b43372730",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AuthRepoImpl{
				db: tt.fields.db,
			}
			got, err := r.SaveUser(tt.args.ctx, tt.args.login, tt.args.passHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthRepoImpl.SaveUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AuthRepoImpl.SaveUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthRepoImpl_User(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rows := mock.
		NewRows([]string{"id", "login", "pass_hash"}).
		AddRow("31487452-31d9-4b1f-a7f8-c00b43372730", "user", []byte{})

	query := `
	SELECT id, login, pass_hash
	FROM usr_users
	WHERE login = .+?
	LIMIT 1;
	`
	mock.ExpectPrepare(query)
	mock.ExpectQuery(query).WillReturnRows(rows)

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx   context.Context
		login string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.UserAuth
		wantErr bool
	}{
		{
			name:   "Success",
			fields: fields{db: sqlxDB},
			args: args{
				ctx:   ctx,
				login: "user",
			},
			want: &models.UserAuth{
				ID:       "31487452-31d9-4b1f-a7f8-c00b43372730",
				Login:    "user",
				PassHash: []byte{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AuthRepoImpl{
				db: tt.fields.db,
			}
			got, err := r.User(tt.args.ctx, tt.args.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthRepoImpl.User() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthRepoImpl.User() = %v, want %v", got, tt.want)
			}
		})
	}
}
