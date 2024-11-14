package auth

import (
	"context"
	"gopasskeeper/internal/lib/jwt"
	"gopasskeeper/internal/logger"
	"gopasskeeper/internal/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func TestAuth_RegisterNewUser(t *testing.T) {
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

	repo := repository.New(sqlxDB)
	log := logger.NewGRPCLogger("account-test")
	secret := "87ccd01b-0b99-4f3a-8422-dd088f22cde0"
	tokenTTL := time.Hour

	type fields struct {
		log         *logger.GRPCLogger
		usrSaver    UserSaver
		usrProvider UserProvider
		secret      string
		tokenTTL    time.Duration
	}
	type args struct {
		ctx      context.Context
		login    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				log:         log,
				usrSaver:    repo.Auth,
				usrProvider: repo.Auth,
				secret:      secret,
				tokenTTL:    tokenTTL,
			},
			args: args{
				ctx:      ctx,
				login:    "user",
				password: "P@ssword!",
			},
			want:    "31487452-31d9-4b1f-a7f8-c00b43372730",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Auth{
				log:         tt.fields.log,
				usrSaver:    tt.fields.usrSaver,
				usrProvider: tt.fields.usrProvider,
				secret:      tt.fields.secret,
				tokenTTL:    tt.fields.tokenTTL,
			}
			got, err := a.RegisterNewUser(tt.args.ctx, tt.args.login, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Auth.RegisterNewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Auth.RegisterNewUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuth_Login(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	password := "P@ssword!"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	rows := mock.
		NewRows([]string{"id", "login", "pass_hash"}).
		AddRow("31487452-31d9-4b1f-a7f8-c00b43372730", "user", passwordHash)

	query := `
	SELECT id, login, pass_hash
	FROM usr_users
	WHERE login = .+?
	LIMIT 1;
	`
	mock.ExpectPrepare(query)
	mock.ExpectQuery(query).WillReturnRows(rows)

	repo := repository.New(sqlxDB)
	log := logger.NewGRPCLogger("account-test")
	secret := "87ccd01b-0b99-4f3a-8422-dd088f22cde0"
	tokenTTL := time.Hour

	jwtManager := jwt.NewJWTManager(secret, tokenTTL)

	type fields struct {
		log         *logger.GRPCLogger
		usrSaver    UserSaver
		usrProvider UserProvider
		secret      string
		tokenTTL    time.Duration
	}
	type args struct {
		ctx      context.Context
		login    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				log:         log,
				usrSaver:    repo.Auth,
				usrProvider: repo.Auth,
				secret:      secret,
				tokenTTL:    tokenTTL,
			},
			args: args{
				ctx:      ctx,
				login:    "user",
				password: password,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Auth{
				log:         tt.fields.log,
				usrSaver:    tt.fields.usrSaver,
				usrProvider: tt.fields.usrProvider,
				secret:      tt.fields.secret,
				tokenTTL:    tt.fields.tokenTTL,
			}
			got, err := a.Login(tt.args.ctx, tt.args.login, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Auth.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if _, err := jwtManager.Verify(got); (err != nil) != tt.wantErr {
				t.Errorf("Auth.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
