package sync

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestSyncRepoImpl_Get(t *testing.T) {
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

	type fields struct {
		db *sqlx.DB
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
			name:   "Success",
			fields: fields{db: sqlxDB},
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
			s := &SyncRepoImpl{
				db: tt.fields.db,
			}
			got, err := s.Get(tt.args.ctx, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("SyncRepoImpl.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SyncRepoImpl.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncRepoImpl_Set(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	query := `
	INSERT INTO syn_timestamps\(uid, timestamp\)
	VALUES \(.+?, now\(\)\)
	ON CONFLICT ON CONSTRAINT syn_timestamps_pk
	DO UPDATE SET timestamp = excluded\.timestamp
	`
	mock.ExpectPrepare(query)
	mock.ExpectQuery(query).WillReturnRows()

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx context.Context
		uid string
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
				ctx: ctx,
				uid: "31487452-31d9-4b1f-a7f8-c00b43372730",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncRepoImpl{
				db: tt.fields.db,
			}
			if err := s.Set(tt.args.ctx, tt.args.uid); (err != nil) != tt.wantErr {
				t.Errorf("SyncRepoImpl.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
