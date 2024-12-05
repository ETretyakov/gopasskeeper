package accounts

import (
	"context"
	"gopasskeeper/internal/config"
	"gopasskeeper/internal/domain/models"
	"gopasskeeper/internal/lib/crypto"
	"gopasskeeper/internal/logger"

	"github.com/pkg/errors"
)

var (
	// ErrAccountNotFound is an error variable to define account not found errors.
	ErrAccountNotFound = errors.New("account hasn't been found")
)

// AccountStorage is an interface to describe account storage methods.
type AccountStorage interface {
	Add(
		ctx context.Context,
		uid string,
		login string,
		server string,
		password string,
		meta string,
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

// SyncStorage is an interface to descrube sync methods.
type SyncStorage interface {
	Set(ctx context.Context, uid string) error
}

// Accounts is a structure to define accounts service.
type Accounts struct {
	log             *logger.GRPCLogger
	fernetEncryptor *crypto.FernetEncryptor
	accountStorage  AccountStorage
	syncStorage     SyncStorage
}

// New is a builder function for accounts service.
func New(
	cfg *config.SecurityConfig,
	accountStorage AccountStorage,
	syncStorage SyncStorage,
) (*Accounts, error) {
	fernetEncryptor, err := crypto.NewFernet(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build fernet encryptor")
	}

	return &Accounts{
		log:             logger.NewGRPCLogger("accounts"),
		fernetEncryptor: fernetEncryptor,
		accountStorage:  accountStorage,
		syncStorage:     syncStorage,
	}, nil
}

// Add is Accounts method to add account record.
func (a *Accounts) Add(
	ctx context.Context,
	uid string,
	login string,
	server string,
	password string,
	meta string,
) (*models.Message, error) {
	const op = "Accounts.Add"
	log := a.log.WithOperator(op)
	log.Info("adding account", "uid", uid)

	encPassword, err := a.fernetEncryptor.Encrypt([]byte(password))
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt password")
	}

	encMeta, err := a.fernetEncryptor.Encrypt([]byte(meta))
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt meta")
	}

	accountID, err := a.accountStorage.Add(
		ctx,
		uid,
		login,
		server,
		string(encPassword[:]),
		string(encMeta[:]),
	)
	if err != nil {
		log.Error("failed to save account", err)
		return nil, errors.Wrap(err, op)
	}

	if err := a.syncStorage.Set(ctx, uid); err != nil {
		log.Error("failed to save account", err)
		return nil, errors.Wrap(err, op)
	}

	return &models.Message{
		Status: true,
		Msg:    "Account added: account id - " + accountID,
	}, nil
}

// GetSecret is Accounts method to get account record with the secret.
func (a *Accounts) GetSecret(
	ctx context.Context,
	uid string,
	accountID string,
) (*models.AccountSecret, error) {
	const op = "Accounts.GetSecret"
	log := a.log.WithOperator(op)
	log.Info("getting secret", "uid", uid)

	accountSecret, err := a.accountStorage.GetSecret(ctx, uid, accountID)
	if err != nil {
		log.Error("failed to get account secret", err)
		return nil, errors.Wrap(err, op)
	}

	decPassword, err := a.fernetEncryptor.Decrypt([]byte(accountSecret.Password))
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt password")
	}
	accountSecret.Password = string(decPassword[:])

	decMeta, err := a.fernetEncryptor.Decrypt([]byte(accountSecret.Meta))
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt meta")
	}
	accountSecret.Meta = string(decMeta[:])

	return accountSecret, nil
}

// Search is Accounts method to search account records.
func (a *Accounts) Search(
	ctx context.Context,
	uid string,
	schema *models.AccountSearchRequest,
) (*models.AccountSearchResponse, error) {
	const op = "Accounts.Search"
	log := a.log.WithOperator(op)
	log.Info("searching items", "uid", uid)

	accountSearchResponse, err := a.accountStorage.Search(ctx, uid, schema)
	if err != nil {
		log.Error("failed to search accounts", err)
		return nil, errors.Wrap(err, op)
	}

	return accountSearchResponse, nil
}

// Remove is Accounts method to remove an account record.
func (a *Accounts) Remove(
	ctx context.Context,
	uid string,
	accountID string,
) (*models.Message, error) {
	const op = "Accounts.Search"
	log := a.log.WithOperator(op)
	log.Info("removing record", "uid", uid)

	if err := a.accountStorage.Remove(ctx, uid, accountID); err != nil {
		log.Error("failed to remove account", err)
		return nil, errors.Wrap(err, op)
	}

	if err := a.syncStorage.Set(ctx, uid); err != nil {
		log.Error("failed to save account", err)
		return nil, errors.Wrap(err, op)
	}

	return &models.Message{
		Status: true,
		Msg:    "Account removed: account id - " + accountID,
	}, nil
}
