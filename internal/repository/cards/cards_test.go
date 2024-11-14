package cards

import (
	"context"
	"gopasskeeper/internal/domain/models"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestCardsRepoImpl_Add(t *testing.T) {
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
	INSERT INTO public\.sec_cards\(uid, name, number, mask, month, year, cvc, pin\)
	VALUES \(.+?, .+?, .+?, .+?, .+?, .+?, .+?, .+?\)
	RETURNING id;
	`
	mock.ExpectPrepare(query)
	mock.ExpectQuery(query).WillReturnRows(rows)

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		uid    string
		name   string
		number string
		mask   string
		month  int32
		year   int32
		cvc    string
		pin    string
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
				ctx:    ctx,
				uid:    "31487452-31d9-4b1f-a7f8-c00b43372730",
				name:   "VISA",
				number: "4242424242424242",
				mask:   "**** **** **** 4242",
				month:  1,
				year:   2025,
				cvc:    "777",
				pin:    "1111",
			},
			want:    "d89b92df-44e8-4d66-857a-bf7ec0a61556",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CardsRepoImpl{
				db: tt.fields.db,
			}
			got, err := r.Add(tt.args.ctx, tt.args.uid, tt.args.name, tt.args.number, tt.args.mask, tt.args.month, tt.args.year, tt.args.cvc, tt.args.pin)
			if (err != nil) != tt.wantErr {
				t.Errorf("CardsRepoImpl.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CardsRepoImpl.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardsRepoImpl_GetSecret(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rows := mock.
		NewRows([]string{"name", "number", "month", "year", "cvc", "pin"}).
		AddRow("VISA", "4242424242424242", 1, 2025, "777", "1111")

	query := `
	SELECT sc\.name   AS \"name\",
		   sc\.number AS \"number\",
		   sc\.month  AS \"month\",
		   sc\.year   AS \"year\",
		   sc\.cvc    AS \"cvc\",
		   sc\.pin    AS \"pin\"
	FROM sec_cards sc
	WHERE sc\.uid = .+? AND 
	      sc\.id  = .+?
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
		cardID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.CardSecret
		wantErr bool
	}{
		{
			name:   "Success",
			fields: fields{db: sqlxDB},
			args: args{
				ctx:    ctx,
				uid:    "31487452-31d9-4b1f-a7f8-c00b43372730",
				cardID: "d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
			want: &models.CardSecret{
				Name:   "VISA",
				Number: "4242424242424242",
				Month:  1,
				Year:   2025,
				CVC:    "777",
				PIN:    "1111",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CardsRepoImpl{
				db: tt.fields.db,
			}
			got, err := r.GetSecret(tt.args.ctx, tt.args.uid, tt.args.cardID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CardsRepoImpl.GetSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CardsRepoImpl.GetSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardsRepoImpl_Search(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rows := mock.
		NewRows([]string{"id", "name", "mask"}).
		AddRow("d89b92df-44e8-4d66-857a-bf7ec0a61556", "VISA", "**** **** **** 4242")

	query := `
	SELECT sc\.id   AS \"id\",
	       sc\.name AS \"name\",
		   sc\.mask AS \"mask\"
	FROM sec_cards sc
	WHERE sc\.uid = .+? AND
		  sc\.name ILIKE .+?
	ORDER BY sc\.name
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
	FROM sec_cards sc
	WHERE sc\.uid = .+? AND
		  sc\.name ILIKE .+?;
	`
	mock.ExpectPrepare(countQuery)
	mock.ExpectQuery(countQuery).WillReturnRows(countRows)

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		uid    string
		schema *models.CardSearchRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.CardSearchResponse
		wantErr bool
	}{
		{
			name:   "Success",
			fields: fields{db: sqlxDB},
			args: args{
				ctx: ctx,
				uid: "31487452-31d9-4b1f-a7f8-c00b43372730",
				schema: &models.CardSearchRequest{
					Substring: "",
					Offset:    0,
					Limit:     100,
				},
			},
			want: &models.CardSearchResponse{
				Count: 1,
				Items: []*models.CardSearchItem{
					{
						ID:   "d89b92df-44e8-4d66-857a-bf7ec0a61556",
						Name: "VISA",
						Mask: "**** **** **** 4242",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CardsRepoImpl{
				db: tt.fields.db,
			}
			got, err := r.Search(tt.args.ctx, tt.args.uid, tt.args.schema)
			if (err != nil) != tt.wantErr {
				t.Errorf("CardsRepoImpl.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CardsRepoImpl.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardsRepoImpl_Remove(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	query := `
	DELETE FROM sec_cards
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
		cardID string
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
				cardID: "d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CardsRepoImpl{
				db: tt.fields.db,
			}
			if err := r.Remove(tt.args.ctx, tt.args.uid, tt.args.cardID); (err != nil) != tt.wantErr {
				t.Errorf("CardsRepoImpl.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
