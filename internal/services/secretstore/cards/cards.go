package cards

import (
	"context"
	"gopasskeeper/internal/config"
	"gopasskeeper/internal/domain/models"
	"gopasskeeper/internal/lib/crypto"
	"gopasskeeper/internal/logger"

	"github.com/pkg/errors"
)

var (
	// ErrCardNotFound is an error variable to define card not found errors.
	ErrCardNotFound = errors.New("card hasn't been found")
)

// CardStorage is an interface to describe card storage methods.
type CardStorage interface {
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

// SyncStorage is an interface to descrube sync methods.
type SyncStorage interface {
	Set(ctx context.Context, uid string) error
}

// Cards is a structure to define cards service.
type Cards struct {
	log             *logger.GRPCLogger
	fernetEncryptor *crypto.FernetEncryptor
	cardStorage     CardStorage
	syncStorage     SyncStorage
}

// New is a builder function for cards service.
func New(
	cfg *config.SecurityConfig,
	cardStorage CardStorage,
	syncStorage SyncStorage,
) (*Cards, error) {
	fernetEncryptor, err := crypto.NewFernet(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build fernet encryptor")
	}

	return &Cards{
		log:             logger.NewGRPCLogger("cards"),
		fernetEncryptor: fernetEncryptor,
		cardStorage:     cardStorage,
		syncStorage:     syncStorage,
	}, nil
}

// Add is Cards method to add card record.
func (c *Cards) Add(
	ctx context.Context,
	uid string,
	name string,
	number string,
	mask string,
	month int32,
	year int32,
	cvc string,
	pin string,
) (*models.Message, error) {
	const op = "Cards.Add"
	log := c.log.WithOperator(op)
	log.Info("adding card", "uid", uid)

	encCVC, err := c.fernetEncryptor.Encrypt([]byte(cvc))
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt cvc")
	}

	encPIN, err := c.fernetEncryptor.Encrypt([]byte(pin))
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt pin")
	}

	cardID, err := c.cardStorage.Add(
		ctx,
		uid,
		name,
		number,
		mask,
		month,
		year,
		string(encCVC[:]),
		string(encPIN[:]),
	)
	if err != nil {
		log.Error("failed to save card", err)
		return nil, errors.Wrap(err, op)
	}

	if err := c.syncStorage.Set(ctx, uid); err != nil {
		log.Error("failed to save card", err)
		return nil, errors.Wrap(err, op)
	}

	return &models.Message{
		Status: true,
		Msg:    "Card added: card id - " + cardID,
	}, nil
}

// GetSecret is Cards method to get card record with the secret.
func (c *Cards) GetSecret(
	ctx context.Context,
	uid string,
	cardID string,
) (*models.CardSecret, error) {
	const op = "Cards.GetSecret"
	log := c.log.WithOperator(op)
	log.Info("getting secret", "uid", uid)

	cardSecret, err := c.cardStorage.GetSecret(ctx, uid, cardID)
	if err != nil {
		log.Error("failed to get card secret", err)
		return nil, errors.Wrap(err, op)
	}

	decCVC, err := c.fernetEncryptor.Decrypt([]byte(cardSecret.CVC))
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt cvc")
	}
	cardSecret.CVC = string(decCVC[:])

	decPIN, err := c.fernetEncryptor.Decrypt([]byte(cardSecret.PIN))
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt pin")
	}
	cardSecret.PIN = string(decPIN[:])

	return cardSecret, nil
}

// Search is Cards method to search card records.
func (c *Cards) Search(
	ctx context.Context,
	uid string,
	schema *models.CardSearchRequest,
) (*models.CardSearchResponse, error) {
	const op = "Cards.Search"
	log := c.log.WithOperator(op)
	log.Info("searching items", "uid", uid)

	cardSearchResponse, err := c.cardStorage.Search(ctx, uid, schema)
	if err != nil {
		log.Error("failed to search cards", err)
		return nil, errors.Wrap(err, op)
	}

	return cardSearchResponse, nil
}

// Remove is Cards method to remove an card record.
func (c *Cards) Remove(
	ctx context.Context,
	uid string,
	cardID string,
) (*models.Message, error) {
	const op = "Cards.Search"
	log := c.log.WithOperator(op)
	log.Info("removing record", "uid", uid)

	if err := c.cardStorage.Remove(ctx, uid, cardID); err != nil {
		log.Error("failed to remove card", err)
		return nil, errors.Wrap(err, op)
	}

	if err := c.syncStorage.Set(ctx, uid); err != nil {
		log.Error("failed to save account", err)
		return nil, errors.Wrap(err, op)
	}

	return &models.Message{
		Status: true,
		Msg:    "Card removed: card id - " + cardID,
	}, nil
}
