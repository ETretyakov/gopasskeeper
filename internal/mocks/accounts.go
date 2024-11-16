package mocks

func (mdb *MockedDB) AccountAddMockedDB(accountID string) *MockedDB {
	rows := mdb.mock.
		NewRows([]string{"id"}).
		AddRow(accountID)

	query := `
	INSERT INTO public\.sec_accounts\(uid, login, password, server\)
	VALUES \(.+?, .+?, .+?, .+?\)
	RETURNING id;
	`
	mdb.mock.ExpectPrepare(query)
	mdb.mock.ExpectQuery(query).WillReturnRows(rows)

	return mdb
}

func (mdb *MockedDB) AccountGetSecretMockedDB(login, server, password string) *MockedDB {
	rows := mdb.mock.
		NewRows([]string{"login", "server", "password"}).
		AddRow(login, server, password)

	query := `
	SELECT sa\.login    AS \"login\",
		   sa\.server   AS \"server\",
		   sa\.password AS \"password\"
	FROM sec_accounts sa
	WHERE sa\.uid = .+? AND 
	      sa\.id  = .+?
	LIMIT 1;
	`
	mdb.mock.ExpectPrepare(query)
	mdb.mock.ExpectQuery(query).WillReturnRows(rows)

	return mdb
}

func (mdb *MockedDB) AccountSearchMockedDB(id, login, server string) *MockedDB {
	rows := mdb.mock.
		NewRows([]string{"id", "login", "server"}).
		AddRow(id, login, server)

	query := `
	SELECT sa\.id       AS \"id\",
	       sa\.login    AS \"login\",
		   sa\.server   AS \"server\"
	FROM sec_accounts sa
	WHERE sa\.uid = .+? AND
		  \(sa\.server ILIKE .+? OR
		   sa\.login  ILIKE .+?\)
	ORDER BY sa\.server, sa\.login
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
	FROM sec_accounts sa
	WHERE sa\.uid = .+? AND
		  \(sa\.server ILIKE .+? OR
		   sa\.login  ILIKE .+?\);
	`
	mdb.mock.ExpectPrepare(countQuery)
	mdb.mock.ExpectQuery(countQuery).WillReturnRows(countRows)

	return mdb
}

func (mdb *MockedDB) AccountRemoveMockedDB() *MockedDB {
	query := `
	DELETE FROM sec_accounts
	WHERE uid = .+? AND
	      id  = .+?;
	`
	mdb.mock.ExpectPrepare(query)
	mdb.mock.ExpectQuery(query).WillReturnRows()

	return mdb
}
