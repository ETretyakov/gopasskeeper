package mocks

func (mdb *MockedDB) CardAddMockedDB(id string) *MockedDB {
	rows := mdb.mock.
		NewRows([]string{"id"}).
		AddRow(id)

	query := `
	INSERT INTO public\.sec_cards\(uid, name, number, mask, month, year, cvc, pin, meta\)
	VALUES \(.+?, .+?, .+?, .+?, .+?, .+?, .+?, .+?, .+?\)
	RETURNING id;
	`
	mdb.mock.ExpectPrepare(query)
	mdb.mock.ExpectQuery(query).WillReturnRows(rows)

	return mdb
}

func (mdb *MockedDB) CardGetSecretMockedDB(
	name, number string,
	month, year int32,
	cvc, pin, meta string,
) *MockedDB {
	rows := mdb.mock.
		NewRows([]string{"name", "number", "month", "year", "cvc", "pin", "meta"}).
		AddRow(name, number, month, year, cvc, pin, meta)

	query := `
	SELECT sc\.name   AS \"name\",
		   sc\.number AS \"number\",
		   sc\.month  AS \"month\",
		   sc\.year   AS \"year\",
		   sc\.cvc    AS \"cvc\",
		   sc\.pin    AS \"pin\",
		   sc\.meta   AS \"meta\"
	FROM sec_cards sc
	WHERE sc\.uid = .+? AND 
	      sc\.id  = .+?
	LIMIT 1;
	`
	mdb.mock.ExpectPrepare(query)
	mdb.mock.ExpectQuery(query).WillReturnRows(rows)

	return mdb
}

func (mdb *MockedDB) CardSearchMockedDB(id, name, mask string) *MockedDB {
	rows := mdb.mock.
		NewRows([]string{"id", "name", "mask"}).
		AddRow(id, name, mask)

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
	mdb.mock.ExpectPrepare(query)
	mdb.mock.ExpectQuery(query).WillReturnRows(rows)

	countRows := mdb.mock.
		NewRows([]string{"count"}).
		AddRow(1)

	countQuery := `
	SELECT count\(\*\) AS \"count\"
	FROM sec_cards sc
	WHERE sc\.uid = .+? AND
		sc\.name ILIKE .+?;
	`
	mdb.mock.ExpectPrepare(countQuery)
	mdb.mock.ExpectQuery(countQuery).WillReturnRows(countRows)

	return mdb
}

func (mdb *MockedDB) CardRemoveMockedDB() *MockedDB {
	query := `
	DELETE FROM sec_cards
	WHERE uid = .+? AND
		  id  = .+?;
	`
	mdb.mock.ExpectPrepare(query)
	mdb.mock.ExpectQuery(query).WillReturnRows()

	return mdb
}
