package mocks

func (mdb *MockedDB) NoteAddMockedDB(id string) *MockedDB {
	rows := mdb.mock.
		NewRows([]string{"id"}).
		AddRow(id)

	query := `
	INSERT INTO public\.sec_notes\(uid, name, content\)
	VALUES \(.+?, .+?, .+?\)
	RETURNING id;
	`
	mdb.mock.ExpectPrepare(query)
	mdb.mock.ExpectQuery(query).WillReturnRows(rows)

	return mdb
}

func (mdb *MockedDB) NoteGetSecretMockedDB(name, content string) *MockedDB {
	rows := mdb.mock.
		NewRows([]string{"name", "content"}).
		AddRow(name, content)

	query := `
	SELECT sn\.name   AS \"name\",
		   sn\.content AS \"content\"
	FROM sec_notes sn
	WHERE sn\.uid = .+? AND 
			sn\.id  = .+?
	LIMIT 1;
	`
	mdb.mock.ExpectPrepare(query)
	mdb.mock.ExpectQuery(query).WillReturnRows(rows)

	return mdb
}

func (mdb *MockedDB) NoteSearchMockedDB(id, name string) *MockedDB {
	rows := mdb.mock.
		NewRows([]string{"id", "name"}).
		AddRow(id, name)

	query := `
	SELECT sn\.id   AS \"id\",
	       sn\.name AS \"name\"
	FROM sec_notes sn
	WHERE sn\.uid = .+? AND
		  sn\.name ILIKE .+?
	ORDER BY sn\.name
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
	FROM sec_notes sn
	WHERE sn\.uid = .+? AND
		sn\.name ILIKE .+?;
	`
	mdb.mock.ExpectPrepare(countQuery)
	mdb.mock.ExpectQuery(countQuery).WillReturnRows(countRows)

	return mdb
}

func (mdb *MockedDB) NoteRemoveMockedDB() *MockedDB {
	query := `
	DELETE FROM sec_notes
	WHERE uid = .+? AND
	      id  = .+?;
	`
	mdb.mock.ExpectPrepare(query)
	mdb.mock.ExpectQuery(query).WillReturnRows()

	return mdb
}
