package api

import (
	"gopasskeeper/internal/clients/secretstore"
	accountsv1 "gopasskeeper/internal/grpc/secretstore/accounts/gen/accounts"
	cardsv1 "gopasskeeper/internal/grpc/secretstore/cards/gen/cards"
	notesv1 "gopasskeeper/internal/grpc/secretstore/notes/gen/notes"

	"sync"
)

type AuthAPI interface {
	SetClient(client *secretstore.Client)

	Login(endpoint, login, password string) error
}

type AccountsAPI interface {
	SetClient(client *secretstore.Client)

	SearchAccount(substring string, offset uint64, limit uint32) (*accountsv1.AccountSearchResponse, error)
	GetAccount(uuid string) (string, error)
	AddAccount(server, login, password string) error
	RemoveAccount(secredID string) error
}

type CardsAPI interface {
	SearchCard(substring string, offset uint64, limit uint32) (*cardsv1.CardSearchResponse, error)
	GetCard(uuid string) (string, error)
	AddCard(name, number string, month, year int32, ccv, pin string) error
	RemoveCard(secredID string) error
}

type NotesAPI interface {
	SearchNote(substring string, offset uint64, limit uint32) (*notesv1.NoteSearchResponse, error)
	GetNote(uuid string) (string, error)
	AddNote(name, content string) error
	RemoveNote(secredID string) error
}

type FilesAPI interface {
	SearchFile(substring string, offset uint64, limit uint32) (*notesv1.NoteSearchResponse, error)
	GetFile(uuid, filePath string) (string, error)
	AddFile(name, filePath string) error
	RemoveNote(secredID string) error
}

type SyncAPI interface {
	Outdated() (bool, error)
}

type APIConfig interface {
	CertPath() string
}

type API struct {
	mu     sync.RWMutex
	client *secretstore.Client
	states *States
	cfg    APIConfig
}

func New(cfg APIConfig) *API {
	return &API{
		mu:     sync.RWMutex{},
		states: NewStates(),
		cfg:    cfg,
	}
}

func (a *API) SetClient(client *secretstore.Client) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.client = client
}
