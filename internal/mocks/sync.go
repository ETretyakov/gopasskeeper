package mocks

func (mdb *MockedDB) AddSyncMocks() *MockedDB {
	syncQuery := `
	INSERT INTO syn_timestamps\(uid, timestamp\)
	VALUES \(.+?, now\(\)\)
	ON CONFLICT ON CONSTRAINT syn_timestamps_pk
	DO UPDATE SET timestamp = excluded\.timestamp
	`
	mdb.mock.ExpectPrepare(syncQuery)
	mdb.mock.ExpectQuery(syncQuery).WillReturnRows()

	return mdb
}
