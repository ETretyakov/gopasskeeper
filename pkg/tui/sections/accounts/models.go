package accounts

type SearchInput struct {
	value  string
	offset uint64
	limit  uint32
	step   int
}

type AccountAdd struct {
	server   string
	login    string
	password string
}
