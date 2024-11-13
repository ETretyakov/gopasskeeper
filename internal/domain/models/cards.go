package models

// CardSecret is a structure to describe card secret.
type CardSecret struct {
	Name   string `db:"name"`
	Number string `db:"number"`
	Month  int32  `db:"month"`
	Year   int32  `db:"year"`
	CVC    string `db:"cvc"`
	PIN    string `db:"pin"`
}

// CardSearchRequest is a structure to describe card search request.
type CardSearchRequest struct {
	Substring string
	Offset    uint64
	Limit     uint32
}

// CardSearchItem is a structure to describe card search item.
type CardSearchItem struct {
	ID   string `db:"id"`
	Name string `db:"name"`
	Mask string `db:"mask"`
}

// CardSearchResponse is a structure to describe account search response.
type CardSearchResponse struct {
	Count uint64
	Items []*CardSearchItem
}
