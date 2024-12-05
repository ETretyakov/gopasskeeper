package models

// FileSecret is a structure to describe card secret.
type FileSecret struct {
	Name    string
	Content []byte
	Meta    string
}

// FileSearchRequest is a structure to describe card search request.
type FileSearchRequest struct {
	Substring string
	Offset    uint64
	Limit     uint32
}

// FileSearchItem is a structure to describe card search item.
type FileSearchItem struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

// FileSearchResponse is a structure to describe account search response.
type FileSearchResponse struct {
	Count uint64
	Items []*FileSearchItem
}
