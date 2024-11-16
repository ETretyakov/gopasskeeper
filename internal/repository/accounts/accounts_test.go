package accounts

import (
	"context"
	"gopasskeeper/internal/domain/models"
	"gopasskeeper/internal/mocks"
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestAccountsRepoImpl_Add(t *testing.T) {
	ctx := context.Background()
	accountID := "d89b92df-44e8-4d66-857a-bf7ec0a61556"
	mockedDB := mocks.NewDB(t).AccountAddMockedDB(accountID)

	type fields struct {
		db *sqlx.DB
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
		want    string
		wantErr bool
	}{
		{
			name:   "Success 1",
			fields: fields{db: mockedDB.Get()},
			args: args{
				ctx:      ctx,
				uid:      "31487452-31d9-4b1f-a7f8-c00b43372730",
				login:    "user",
				server:   "https://test.com",
				password: "P@ssWord!",
			},
			want:    accountID,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AccountsRepoImpl{
				db: tt.fields.db,
			}
			got, err := r.Add(tt.args.ctx, tt.args.uid, tt.args.login, tt.args.server, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountsRepoImpl.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AccountsRepoImpl.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccountsRepoImpl_GetSecret(t *testing.T) {
	ctx := context.Background()

	const (
		login    = "user"
		server   = "https://test.com"
		password = "P@ssWord!"
	)

	mockedDB := mocks.NewDB(t).
		AccountGetSecretMockedDB(login, server, password)

	type fields struct {
		db *sqlx.DB
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
			name:   "Success",
			fields: fields{db: mockedDB.Get()},
			args: args{
				ctx:       ctx,
				uid:       "31487452-31d9-4b1f-a7f8-c00b43372730",
				accountID: "d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
			want: &models.AccountSecret{
				Login:    login,
				Server:   server,
				Password: password,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AccountsRepoImpl{
				db: tt.fields.db,
			}
			got, err := r.GetSecret(tt.args.ctx, tt.args.uid, tt.args.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountsRepoImpl.GetSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AccountsRepoImpl.GetSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccountsRepoImpl_Search(t *testing.T) {
	ctx := context.Background()

	const (
		id     = "d89b92df-44e8-4d66-857a-bf7ec0a61556"
		login  = "user"
		server = "https://test.com"
	)

	mockedDB := mocks.NewDB(t).
		AccountSearchMockedDB(id, login, server)

	type fields struct {
		db *sqlx.DB
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
			name:   "Success",
			fields: fields{db: mockedDB.Get()},
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
			r := &AccountsRepoImpl{
				db: tt.fields.db,
			}
			got, err := r.Search(tt.args.ctx, tt.args.uid, tt.args.schema)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountsRepoImpl.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AccountsRepoImpl.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccountsRepoImpl_Remove(t *testing.T) {
	ctx := context.Background()

	mockedDB := mocks.NewDB(t).
		AccountRemoveMockedDB()

	type fields struct {
		db *sqlx.DB
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
		wantErr bool
	}{
		{
			name:   "Success",
			fields: fields{db: mockedDB.Get()},
			args: args{
				ctx:       ctx,
				uid:       "31487452-31d9-4b1f-a7f8-c00b43372730",
				accountID: "d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &AccountsRepoImpl{
				db: tt.fields.db,
			}
			if err := r.Remove(tt.args.ctx, tt.args.uid, tt.args.accountID); (err != nil) != tt.wantErr {
				t.Errorf("AccountsRepoImpl.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
