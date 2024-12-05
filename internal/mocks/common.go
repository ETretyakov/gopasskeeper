package mocks

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

type MockedDB struct {
	db   *sqlx.DB
	mock sqlmock.Sqlmock
}

func NewDB(t *testing.T) *MockedDB {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	return &MockedDB{
		db:   sqlxDB,
		mock: mock,
	}
}

func (mdb *MockedDB) Get() *sqlx.DB {
	return mdb.db
}
