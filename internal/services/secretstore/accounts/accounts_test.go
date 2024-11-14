package accounts

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

func TestAccounts_Add(t *testing.T) {
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
	INSERT INTO public\.sec_accounts\(uid, login, password, server\)
	VALUES \(.+?, .+?, .+?, .+?\)
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
		accountStorage  AccountStorage
		syncStorage     SyncStorage
	}
	type args struct {
		ctx      context.Context
		uid      string
		login    string
		server   string
		password string
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
				accountStorage:  repo.Accounts,
				syncStorage:     repo.Sync,
			},
			args: args{
				ctx:      ctx,
				uid:      "31487452-31d9-4b1f-a7f8-c00b43372730",
				login:    "user",
				server:   "https://server.test",
				password: "P@ssword!",
			},
			want: &models.Message{
				Status: true,
				Msg:    "Account added: account id - d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Accounts{
				log:             tt.fields.log,
				fernetEncryptor: tt.fields.fernetEncryptor,
				accountStorage:  tt.fields.accountStorage,
				syncStorage:     tt.fields.syncStorage,
			}
			got, err := a.Add(tt.args.ctx, tt.args.uid, tt.args.login, tt.args.server, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Accounts.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Accounts.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccounts_GetSecret(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	fernetEncryptor, err := crypto.NewFernet(
		&config.SecurityConfig{
			Fernet: "QijSv1fl9KAz733U_Rjxc2ribjQpJguYP2C5ezrQcwA=",
		},
	)

	password := "P@ssWord!"
	encPass, _ := fernetEncryptor.Encrypt([]byte(password))
	rows := mock.
		NewRows([]string{"login", "server", "password"}).
		AddRow("user", "https://test.server", string(encPass[:]))

	query := `
	SELECT sa\.login    AS \"login\",
		   sa\.server   AS \"server\",
		   sa\.password AS \"password\"
	FROM sec_accounts sa
	WHERE sa\.uid = .+? AND 
	      sa\.id  = .+?
	LIMIT 1;
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

	if err != nil {
		return
	}

	type fields struct {
		log             *logger.GRPCLogger
		fernetEncryptor *crypto.FernetEncryptor
		accountStorage  AccountStorage
		syncStorage     SyncStorage
	}
	type args struct {
		ctx       context.Context
		uid       string
		accountID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.AccountSecret
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				log:             log,
				fernetEncryptor: fernetEncryptor,
				accountStorage:  repo.Accounts,
				syncStorage:     repo.Sync,
			},
			args: args{
				ctx:       ctx,
				uid:       "31487452-31d9-4b1f-a7f8-c00b43372730",
				accountID: "d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
			want: &models.AccountSecret{
				Login:    "user",
				Server:   "https://test.server",
				Password: password,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Accounts{
				log:             tt.fields.log,
				fernetEncryptor: tt.fields.fernetEncryptor,
				accountStorage:  tt.fields.accountStorage,
				syncStorage:     tt.fields.syncStorage,
			}
			got, err := a.GetSecret(tt.args.ctx, tt.args.uid, tt.args.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Accounts.GetSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Accounts.GetSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccounts_Search(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rows := mock.
		NewRows([]string{"id", "login", "server"}).
		AddRow("d89b92df-44e8-4d66-857a-bf7ec0a61556", "user", "https://test.com")

	query := `
	SELECT sa\.id       AS \"id\",
	       sa\.login    AS \"login\",
		   sa\.server   AS \"server\"
	FROM sec_accounts sa
	WHERE sa\.uid = .+? AND
		  \(sa\.server ILIKE .+? OR
		   sa\.login  ILIKE .+?\)
	ORDER BY sa\.server, sa\.login
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
	FROM sec_accounts sa
	WHERE sa\.uid = .+? AND
		  \(sa\.server ILIKE .+? OR
		   sa\.login  ILIKE .+?\);
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
		accountStorage  AccountStorage
		syncStorage     SyncStorage
	}
	type args struct {
		ctx    context.Context
		uid    string
		schema *models.AccountSearchRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.AccountSearchResponse
		wantErr bool
	}{
		{
			fields: fields{
				log:             log,
				fernetEncryptor: fernetEncryptor,
				accountStorage:  repo.Accounts,
				syncStorage:     repo.Sync,
			},
			args: args{
				ctx: ctx,
				uid: "31487452-31d9-4b1f-a7f8-c00b43372730",
				schema: &models.AccountSearchRequest{
					Substring: "",
					Offset:    0,
					Limit:     100,
				},
			},
			want: &models.AccountSearchResponse{
				Count: 1,
				Items: []*models.AccountSearchItem{
					{
						ID:     "d89b92df-44e8-4d66-857a-bf7ec0a61556",
						Login:  "user",
						Server: "https://test.com",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Accounts{
				log:             tt.fields.log,
				fernetEncryptor: tt.fields.fernetEncryptor,
				accountStorage:  tt.fields.accountStorage,
				syncStorage:     tt.fields.syncStorage,
			}
			got, err := a.Search(tt.args.ctx, tt.args.uid, tt.args.schema)
			if (err != nil) != tt.wantErr {
				t.Errorf("Accounts.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Accounts.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccounts_Remove(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	query := `
	DELETE FROM sec_accounts
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
		accountStorage  AccountStorage
		syncStorage     SyncStorage
	}
	type args struct {
		ctx       context.Context
		uid       string
		accountID string
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
				accountStorage:  repo.Accounts,
				syncStorage:     repo.Sync,
			},
			args: args{
				ctx:       ctx,
				uid:       "31487452-31d9-4b1f-a7f8-c00b43372730",
				accountID: "d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
			want: &models.Message{
				Status: true,
				Msg:    "Account removed: account id - d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Accounts{
				log:             tt.fields.log,
				fernetEncryptor: tt.fields.fernetEncryptor,
				accountStorage:  tt.fields.accountStorage,
				syncStorage:     tt.fields.syncStorage,
			}
			got, err := a.Remove(tt.args.ctx, tt.args.uid, tt.args.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Accounts.Remove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Accounts.Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}
