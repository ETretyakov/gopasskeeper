package models

// UserAuth is a structure to represent user object for authentication.
type UserAuth struct {
	ID       string `db:"id"`
	Login    string `db:"login"`
	PassHash []byte `db:"pass_hash"`
}
