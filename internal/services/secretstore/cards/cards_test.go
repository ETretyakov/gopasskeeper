package cards

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

func TestCards_Add(t *testing.T) {
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
		cardStorage     CardStorage
		syncStorage     SyncStorage
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
		want    *models.Message
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				log:             log,
				fernetEncryptor: fernetEncryptor,
				cardStorage:     repo.Cards,
				syncStorage:     repo.Sync,
			},
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
			want: &models.Message{
				Status: true,
				Msg:    "Card added: card id - d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cards{
				log:             tt.fields.log,
				fernetEncryptor: tt.fields.fernetEncryptor,
				cardStorage:     tt.fields.cardStorage,
				syncStorage:     tt.fields.syncStorage,
			}
			got, err := c.Add(tt.args.ctx, tt.args.uid, tt.args.name, tt.args.number, tt.args.mask, tt.args.month, tt.args.year, tt.args.cvc, tt.args.pin)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cards.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cards.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCards_GetSecret(t *testing.T) {
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

	cvc := "777"
	encCVC, _ := fernetEncryptor.Encrypt([]byte(cvc))

	pin := "1111"
	encPIN, _ := fernetEncryptor.Encrypt([]byte(pin))

	rows := mock.
		NewRows([]string{"name", "number", "month", "year", "cvc", "pin"}).
		AddRow("VISA", "4242424242424242", 1, 2025, encCVC, encPIN)

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

	repo := repository.New(sqlxDB)
	log := logger.NewGRPCLogger("accounts-test")

	type fields struct {
		log             *logger.GRPCLogger
		fernetEncryptor *crypto.FernetEncryptor
		cardStorage     CardStorage
		syncStorage     SyncStorage
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
			name: "Success",
			fields: fields{
				log:             log,
				fernetEncryptor: fernetEncryptor,
				cardStorage:     repo.Cards,
				syncStorage:     repo.Sync,
			},
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
				CVC:    cvc,
				PIN:    pin,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cards{
				log:             tt.fields.log,
				fernetEncryptor: tt.fields.fernetEncryptor,
				cardStorage:     tt.fields.cardStorage,
				syncStorage:     tt.fields.syncStorage,
			}
			got, err := c.GetSecret(tt.args.ctx, tt.args.uid, tt.args.cardID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cards.GetSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cards.GetSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCards_Search(t *testing.T) {
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
		cardStorage     CardStorage
		syncStorage     SyncStorage
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
			fields: fields{
				log:             log,
				fernetEncryptor: fernetEncryptor,
				cardStorage:     repo.Cards,
				syncStorage:     repo.Sync,
			},
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
			c := &Cards{
				log:             tt.fields.log,
				fernetEncryptor: tt.fields.fernetEncryptor,
				cardStorage:     tt.fields.cardStorage,
				syncStorage:     tt.fields.syncStorage,
			}
			got, err := c.Search(tt.args.ctx, tt.args.uid, tt.args.schema)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cards.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cards.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCards_Remove(t *testing.T) {
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
		cardStorage     CardStorage
		syncStorage     SyncStorage
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
		want    *models.Message
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				log:             log,
				fernetEncryptor: fernetEncryptor,
				cardStorage:     repo.Cards,
				syncStorage:     repo.Sync,
			},
			args: args{
				ctx:    ctx,
				uid:    "31487452-31d9-4b1f-a7f8-c00b43372730",
				cardID: "d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
			want: &models.Message{
				Status: true,
				Msg:    "Card removed: card id - d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cards{
				log:             tt.fields.log,
				fernetEncryptor: tt.fields.fernetEncryptor,
				cardStorage:     tt.fields.cardStorage,
				syncStorage:     tt.fields.syncStorage,
			}
			got, err := c.Remove(tt.args.ctx, tt.args.uid, tt.args.cardID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cards.Remove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cards.Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}
