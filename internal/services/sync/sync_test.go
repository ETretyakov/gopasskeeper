package sync

import (
	"context"
	"gopasskeeper/internal/logger"
	"gopasskeeper/internal/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestSync_Get(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	rows := mock.
		NewRows([]string{"timestamp"}).
		AddRow("2024-11-14T11:18:03.272454Z")

	query := `
	SELECT st\.timestamp AS \"timestamp\"
	FROM syn_timestamps st
	WHERE st\.uid = .+?
	LIMIT 1;`

	mock.ExpectPrepare(query)
	mock.ExpectQuery(query).WillReturnRows(rows)

	repo := repository.New(sqlxDB)
	log := logger.NewGRPCLogger("sync-test")

	type fields struct {
		log         *logger.GRPCLogger
		syncStorage SyncStorage
	}
	type args struct {
		ctx context.Context
		uid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			fields: fields{
				log:         log,
				syncStorage: repo.Sync,
			},
			args: args{
				ctx: ctx,
				uid: "31487452-31d9-4b1f-a7f8-c00b43372730",
			},
			want:    "2024-11-14T11:18:03.272454Z",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sync{
				log:         tt.fields.log,
				syncStorage: tt.fields.syncStorage,
			}
			got, err := s.Get(tt.args.ctx, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sync.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Sync.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
