package accounts

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

func TestAccounts_Add(t *testing.T) {
	ctx := context.Background()
	accountID := "d89b92df-44e8-4d66-857a-bf7ec0a61556"
	mockedDB := mocks.NewDB(t).
		AccountAddMockedDB(accountID).
		AddSyncMocks()

	repo := repository.New(mockedDB.Get())
	log := logger.NewGRPCLogger("accounts-test")
	fernetEncryptor := mocks.NewFernet(t)

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
		meta     string
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
				Msg:    "Account added: account id - " + accountID,
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
			got, err := a.Add(tt.args.ctx, tt.args.uid, tt.args.login, tt.args.server, tt.args.password, tt.args.meta)
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

	const (
		login    = "user"
		server   = "https://test.com"
		password = "P@ssWord!"
		meta     = "meta"
	)

	fernetEncryptor := mocks.NewFernet(t)
	encPass, _ := fernetEncryptor.Encrypt([]byte(password))
	encMeta, _ := fernetEncryptor.Encrypt([]byte(meta))

	mockedDB := mocks.NewDB(t).
		AccountGetSecretMockedDB(login, server, string(encPass[:]), string(encMeta[:])).
		AddSyncMocks()

	repo := repository.New(mockedDB.Get())
	log := logger.NewGRPCLogger("accounts-test")

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
				Login:    login,
				Server:   server,
				Password: password,
				Meta:     meta,
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

	const (
		id     = "d89b92df-44e8-4d66-857a-bf7ec0a61556"
		login  = "user"
		server = "https://test.com"
	)

	mockedDB := mocks.NewDB(t).
		AccountSearchMockedDB(id, login, server)

	repo := repository.New(mockedDB.Get())
	log := logger.NewGRPCLogger("accounts-test")
	fernetEncryptor := mocks.NewFernet(t)

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
						ID:     id,
						Login:  login,
						Server: server,
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

	mockedDB := mocks.NewDB(t).
		AccountRemoveMockedDB().
		AddSyncMocks()

	repo := repository.New(mockedDB.Get())
	log := logger.NewGRPCLogger("accounts-test")
	fernetEncryptor := mocks.NewFernet(t)

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
