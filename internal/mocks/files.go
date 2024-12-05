package mocks

func (mdb *MockedDB) FileAddMockedDB(id string) *MockedDB {
	rows := mdb.mock.
		NewRows([]string{"id"}).
		AddRow(id)

	query := `
	INSERT INTO public\.sec_files\(uid, name, meta\)
	VALUES \(.+?, .+?, .+?\)
	RETURNING id;
	`
	mdb.mock.ExpectPrepare(query)
	mdb.mock.ExpectQuery(query).WillReturnRows(rows)

	return mdb
}

func (mdb *MockedDB) FileGetSecretMockedDB(name, meta string) *MockedDB {
	rows := mdb.mock.
		NewRows([]string{"name", "meta"}).
		AddRow(name, meta)

	query := `
	SELECT sf\.name   AS \"name\",
	       sf\.meta   AS \"meta\"
	FROM sec_files sf
	WHERE sf\.uid = .+? AND 
	      sf\.id  = .+?
	LIMIT 1;
	`
	mdb.mock.ExpectPrepare(query)
	mdb.mock.ExpectQuery(query).WillReturnRows(rows)

	return mdb
}

func (mdb *MockedDB) FileSearchMockedDB(id, name string) *MockedDB {
	rows := mdb.mock.
		NewRows([]string{"id", "name"}).
		AddRow(id, name)

	query := `
	SELECT sf\.id   AS \"id\",
	       sf\.name AS \"name\"
	FROM sec_files sf
	WHERE sf\.uid = .+? AND
		  sf\.name ILIKE .+?
	ORDER BY sf\.name
	OFFSET .+?
	LIMIT  .+?;
	`
	mdb.mock.ExpectPrepare(query)
	mdb.mock.ExpectQuery(query).WillReturnRows(rows)

	countRows := mdb.mock.
		NewRows([]string{"count"}).
		AddRow(1)

	countQuery := `
	SELECT count\(\*\) AS \"count\"
	FROM sec_files sf
	WHERE sf\.uid = .+? AND
		  sf\.name ILIKE .+?;
	`
	mdb.mock.ExpectPrepare(countQuery)
	mdb.mock.ExpectQuery(countQuery).WillReturnRows(countRows)

	return mdb
}

func (mdb *MockedDB) FileRemoveMockedDB() *MockedDB {
	query := `
	DELETE FROM sec_files
	WHERE uid = .+? AND
	      id  = .+?;
	`
	mdb.mock.ExpectPrepare(query)
	mdb.mock.ExpectQuery(query).WillReturnRows()

	return mdb
}
