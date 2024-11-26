package cards

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

func TestCards_Add(t *testing.T) {
	ctx := context.Background()

	const id = "d89b92df-44e8-4d66-857a-bf7ec0a61556"

	mockedDB := mocks.NewDB(t).
		CardAddMockedDB(id).
		AddSyncMocks()

	repo := repository.New(mockedDB.Get())
	log := logger.NewGRPCLogger("accounts-test")
	fernetEncryptor := mocks.NewFernet(t)

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
		meta   string
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
				meta:   "meta",
			},
			want: &models.Message{
				Status: true,
				Msg:    "Card added: card id - " + id,
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
			got, err := c.Add(tt.args.ctx, tt.args.uid, tt.args.name, tt.args.number, tt.args.mask, tt.args.month, tt.args.year, tt.args.cvc, tt.args.pin, tt.args.meta)
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

	const (
		name   = "VISA"
		number = "4242424242424242"
		month  = 1
		year   = 2025
		cvc    = "777"
		pin    = "1111"
	)

	fernetEncryptor := mocks.NewFernet(t)
	encCVC, _ := fernetEncryptor.Encrypt([]byte(cvc))
	encPIN, _ := fernetEncryptor.Encrypt([]byte(pin))

	mockedDB := mocks.NewDB(t).
		CardGetSecretMockedDB(
			name,
			number,
			month, year,
			string(encCVC[:]), string(encPIN[:]),
		)

	repo := repository.New(mockedDB.Get())
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

	const (
		id   = "d89b92df-44e8-4d66-857a-bf7ec0a61556"
		name = "VISA"
		mask = "**** **** **** 4242"
	)

	mockedDB := mocks.NewDB(t).
		CardSearchMockedDB(id, name, mask)

	repo := repository.New(mockedDB.Get())
	log := logger.NewGRPCLogger("accounts-test")
	fernetEncryptor := mocks.NewFernet(t)

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
						ID:   id,
						Name: name,
						Mask: mask,
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

	mockedDB := mocks.NewDB(t).
		CardRemoveMockedDB().
		AddSyncMocks()

	repo := repository.New(mockedDB.Get())
	log := logger.NewGRPCLogger("accounts-test")
	fernetEncryptor := mocks.NewFernet(t)

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
