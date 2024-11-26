package notes

type SearchInput struct {
	value  string
	offset uint64
	limit  uint32
	step   int
}

type NoteAdd struct {
	name    string
	content string
	meta    string
}
