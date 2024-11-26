package cards

import (
	"context"
	"gopasskeeper/internal/domain/models"
	"gopasskeeper/internal/mocks"
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestCardsRepoImpl_Add(t *testing.T) {
	ctx := context.Background()

	const id = "d89b92df-44e8-4d66-857a-bf7ec0a61556"
	mockedDB := mocks.NewDB(t).CardAddMockedDB(id)

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
		meta   string
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
			fields: fields{db: mockedDB.Get()},
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
			want:    id,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CardsRepoImpl{
				db: tt.fields.db,
			}
			got, err := r.Add(tt.args.ctx, tt.args.uid, tt.args.name, tt.args.number, tt.args.mask, tt.args.month, tt.args.year, tt.args.cvc, tt.args.pin, tt.args.meta)
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

	const (
		name   = "VISA"
		number = "4242424242424242"
		month  = 1
		year   = 2025
		cvc    = "777"
		pin    = "1111"
	)

	mockedDB := mocks.NewDB(t).
		CardGetSecretMockedDB(
			name,
			number,
			month,
			year,
			cvc,
			pin,
		)

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
			fields: fields{db: mockedDB.Get()},
			args: args{
				ctx:    ctx,
				uid:    "31487452-31d9-4b1f-a7f8-c00b43372730",
				cardID: "d89b92df-44e8-4d66-857a-bf7ec0a61556",
			},
			want: &models.CardSecret{
				Name:   name,
				Number: number,
				Month:  month,
				Year:   year,
				CVC:    cvc,
				PIN:    pin,
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

	const (
		id   = "d89b92df-44e8-4d66-857a-bf7ec0a61556"
		name = "VISA"
		mask = "**** **** **** 4242"
	)

	mockedDB := mocks.NewDB(t).
		CardSearchMockedDB(id, name, mask)

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
			fields: fields{db: mockedDB.Get()},
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
	mockedDB := mocks.NewDB(t).
		CardRemoveMockedDB()

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
			fields: fields{db: mockedDB.Get()},
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
