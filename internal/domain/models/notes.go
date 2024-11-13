package models

// NoteSecret is a structure to describe card secret.
type NoteSecret struct {
	Name    string `db:"name"`
	Content string `db:"content"`
}

// NoteSearchRequest is a structure to describe card search request.
type NoteSearchRequest struct {
	Substring string
	Offset    uint64
	Limit     uint32
}

// NoteSearchItem is a structure to describe card search item.
type NoteSearchItem struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

// NoteSearchResponse is a structure to describe account search response.
type NoteSearchResponse struct {
	Count uint64
	Items []*NoteSearchItem
}
