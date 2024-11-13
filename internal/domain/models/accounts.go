package models

// AccountSecret is a structure to describe account secret.
type AccountSecret struct {
	Login    string `db:"login"`
	Server   string `db:"server"`
	Password string `db:"password"`
}

// AccountSearchRequest is a structure to describe account search request.
type AccountSearchRequest struct {
	Substring string
	Offset    uint64
	Limit     uint32
}

// AccountSearchItem is a structure to describe account search item.
type AccountSearchItem struct {
	ID     string `db:"id"`
	Login  string `db:"login"`
	Server string `db:"server"`
}

// AccountSearchResponse is a structure to describe account search response.
type AccountSearchResponse struct {
	Count uint64
	Items []*AccountSearchItem
}
