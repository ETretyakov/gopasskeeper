package cards

type SearchInput struct {
	value  string
	offset uint64
	limit  uint32
	step   int
}

type CardAdd struct {
	name   string
	number string
	month  int32
	year   int32
	cvc    string
	pin    string
}
