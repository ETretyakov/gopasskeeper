package accounts

import (
	"context"
	"gopasskeeper/internal/domain/models"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestAccountsRepoImpl_Add(t *testing.T) {
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
			fields: fields{db: sqlxDB},
			args: args{
				ctx:      ctx,
				uid:      "31487452-31d9-4b1f-a7f8-c00b43372730",
				login:    "user",
				server:   "https://test.com",
				password: "P@ssWord!",
			},
			want:    "d89b92df-44e8-4d66-857a-bf7ec0a61556",
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
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rows := mock.
		NewRows([]string{"login", "server", "password"}).
		AddRow("user", "https://test.com", "P@ssWord!")

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
			fields: fields{db: sqlxDB},
			args: args{
				ctx:       ctx,
				uid:       "31487452-31d9-4b1f-a7f8-c00b43372730",
				accountID: "d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
			want: &models.AccountSecret{
				Login:    "user",
				Server:   "https://test.com",
				Password: "P@ssWord!",
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
			fields: fields{db: sqlxDB},
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
			fields: fields{db: sqlxDB},
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
