package repository

import (
	"context"
	"gopasskeeper/internal/domain/models"
	"gopasskeeper/internal/repository/accounts"
	"gopasskeeper/internal/repository/auth"
	"gopasskeeper/internal/repository/cards"
	"gopasskeeper/internal/repository/files"
	"gopasskeeper/internal/repository/notes"
	"gopasskeeper/internal/repository/sync"

	"github.com/jmoiron/sqlx"
)

// AuthRepo is an interface to declare auth repository methods.
type AuthRepo interface {
	User(ctx context.Context, login string) (*models.UserAuth, error)
	SaveUser(ctx context.Context, login string, passHash []byte) (uid string, err error)
}

// AccountsRepo is an interface to declare accounts repository methods.
type AccountsRepo interface {
	Add(
		ctx context.Context,
		uid string,
		login string,
		server string,
		password string,
	) (string, error)
	GetSecret(
		ctx context.Context,
		uid string,
		accountID string,
	) (*models.AccountSecret, error)
	Search(
		ctx context.Context,
		uid string,
		schema *models.AccountSearchRequest,
	) (*models.AccountSearchResponse, error)
	Remove(
		ctx context.Context,
		uid string,
		accountID string,
	) error
}

// CardRepo is an interface to describe cards storage methods.
type CardRepo interface {
	Add(
		ctx context.Context,
		uid string,
		name string,
		number string,
		mask string,
		month int32,
		year int32,
		cvc string,
		pin string,
	) (string, error)
	GetSecret(
		ctx context.Context,
		uid string,
		cardID string,
	) (*models.CardSecret, error)
	Search(
		ctx context.Context,
		uid string,
		schema *models.CardSearchRequest,
	) (*models.CardSearchResponse, error)
	Remove(
		ctx context.Context,
		uid string,
		cardID string,
	) error
}

// NoteRepo is an interface to describe notes storage methods.
type NoteRepo interface {
	Add(
		ctx context.Context,
		uid string,
		name string,
		content string,
	) (string, error)
	GetSecret(
		ctx context.Context,
		uid string,
		noteID string,
	) (*models.NoteSecret, error)
	Search(
		ctx context.Context,
		uid string,
		schema *models.NoteSearchRequest,
	) (*models.NoteSearchResponse, error)
	Remove(
		ctx context.Context,
		uid string,
		noteID string,
	) error
}

// FileRepo is an interface to describe file storage methods.
type FileRepo interface {
	Add(
		ctx context.Context,
		uid string,
		name string,
	) (string, error)
	GetSecret(
		ctx context.Context,
		uid string,
		fileID string,
	) (string, error)
	Search(
		ctx context.Context,
		uid string,
		schema *models.FileSearchRequest,
	) (*models.FileSearchResponse, error)
	Remove(
		ctx context.Context,
		uid string,
		fileID string,
	) error
}

// SyncRepo is an interface to describe sync storage methods.
type SyncRepo interface {
	Get(ctx context.Context, uid string) (string, error)
	Set(ctx context.Context, uid string) error
}

// Repository is a structure to aggregate all repositories.
type Repository struct {
	Auth     AuthRepo
	Accounts AccountsRepo
	Cards    CardRepo
	Notes    NoteRepo
	Files    FileRepo
	Sync     SyncRepo
}

// New is a builder function for Repository.
func New(db *sqlx.DB) *Repository {
	return &Repository{
		Auth:     auth.New(db),
		Accounts: accounts.New(db),
		Cards:    cards.New(db),
		Notes:    notes.New(db),
		Files:    files.New(db),
		Sync:     sync.New(db),
	}
}
