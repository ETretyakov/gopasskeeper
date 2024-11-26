package files

type SearchInput struct {
	value  string
	offset uint64
	limit  uint32
	step   int
}

type FileAdd struct {
	name     string
	filePath string
	meta     string
}
