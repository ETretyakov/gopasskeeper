package cards

import (
	"context"
	"gopasskeeper/internal/domain/models"
	"gopasskeeper/internal/logger"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// CardsRepoImpl is a structure to implement cards repository.
type CardsRepoImpl struct {
	db *sqlx.DB
}

// New is a builder function for CardsRepoImpl.
func New(db *sqlx.DB) *CardsRepoImpl {
	return &CardsRepoImpl{db: db}
}

// Add is a repository method to add record.
func (r *CardsRepoImpl) Add(
	ctx context.Context,
	uid string,
	name string,
	number string,
	mask string,
	month int32,
	year int32,
	cvc string,
	pin string,
) (string, error) {
	const op = "repository.Cards.Add"

	stmt := `
	INSERT INTO public.sec_cards(uid, name, number, mask, month, year, cvc, pin)
	VALUES (:uid, :name, :number, :mask, :month, :year, :cvc, :pin)
	RETURNING id;
	`

	namedStmt, err := r.db.PrepareNamedContext(ctx, stmt)
	if err != nil {
		logger.Error("failed to prepare named context", err)
		return "", errors.Wrap(err, op)
	}

	arg := map[string]interface{}{
		"uid":    uid,
		"name":   name,
		"number": number,
		"mask":   mask,
		"month":  month,
		"year":   year,
		"cvc":    cvc,
		"pin":    pin,
	}

	var cardID string
	if err := namedStmt.QueryRowxContext(ctx, arg).Scan(&cardID); err != nil {
		logger.Error("failed to query row context", err)
		return "", errors.Wrap(err, op)
	}

	return cardID, nil
}

// GetSecret is a repository method to get record secret.
func (r *CardsRepoImpl) GetSecret(
	ctx context.Context,
	uid string,
	cardID string,
) (*models.CardSecret, error) {
	const op = "repository.Cards.GetSecret"

	stmt := `
	SELECT sc.name   AS "name",
		   sc.number AS "number",
		   sc.month  AS "month",
		   sc.year   AS "year",
		   sc.cvc    AS "cvc",
		   sc.pin    AS "pin"
	FROM sec_cards sc
	WHERE sc.uid = :uid AND 
	      sc.id  = :card_id
	LIMIT 1;
	`

	namedStmt, err := r.db.PrepareNamedContext(ctx, stmt)
	if err != nil {
		logger.Error("failed to prepare named context", err)
		return nil, errors.Wrap(err, op)
	}

	arg := map[string]interface{}{
		"uid":     uid,
		"card_id": cardID,
	}

	var cardSecret models.CardSecret
	if err := namedStmt.QueryRowxContext(ctx, arg).StructScan(&cardSecret); err != nil {
		logger.Error("failed to query row context", err)
		return nil, errors.Wrap(err, op)
	}

	return &cardSecret, nil
}

// Search is a repository method to search records.
func (r *CardsRepoImpl) Search(
	ctx context.Context,
	uid string,
	schema *models.CardSearchRequest,
) (*models.CardSearchResponse, error) {
	const op = "repository.Cards.Search"

	// selecting items
	stmt := `
	SELECT sc.id   AS "id",
	       sc.name AS "name",
		   sc.mask AS "mask"
	FROM sec_cards sc
	WHERE sc.uid = :uid AND
		  sc.name ILIKE :substring
	ORDER BY sc.name
	OFFSET :offset
	LIMIT  :limit;
	`

	namedStmt, err := r.db.PrepareNamedContext(ctx, stmt)
	if err != nil {
		logger.Error("failed to prepare named context", err)
		return nil, errors.Wrap(err, op)
	}

	arg := map[string]interface{}{
		"uid":       uid,
		"substring": `%` + schema.Substring + `%`,
		"offset":    schema.Offset,
		"limit":     schema.Limit,
	}

	rows, err := namedStmt.QueryxContext(ctx, arg)
	if err != nil {
		logger.Error("failed to query row context", err)
		return nil, errors.Wrap(err, op)
	}
	defer rows.Close()

	var items []*models.CardSearchItem
	for rows.Next() {
		var cardSearchItem models.CardSearchItem
		if err := rows.StructScan(&cardSearchItem); err != nil {
			logger.Error("failed to query row context", err)
			return nil, errors.Wrap(err, op)
		}

		items = append(items, &cardSearchItem)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, op)
	}

	// selecting count
	stmt = `
	SELECT count(*) AS "count"
	FROM sec_cards sc
	WHERE sc.uid = :uid AND
		  sc.name ILIKE :substring;
	`

	namedStmt, err = r.db.PrepareNamedContext(ctx, stmt)
	if err != nil {
		logger.Error("failed to prepare named context", err)
		return nil, errors.Wrap(err, op)
	}

	var count uint64
	if err := namedStmt.QueryRowxContext(ctx, arg).Scan(&count); err != nil {
		logger.Error("failed to query row context", err)
		return nil, errors.Wrap(err, op)
	}

	return &models.CardSearchResponse{
		Count: count,
		Items: items,
	}, nil
}

// Remove is a repository method to remove a record.
func (r *CardsRepoImpl) Remove(
	ctx context.Context,
	uid string,
	cardID string,
) error {
	const op = "repository.Cards.Remove"

	stmt := `
	DELETE FROM sec_cards
	WHERE uid = :uid AND
	      id  = :card_id;
	`

	namedStmt, err := r.db.PrepareNamedContext(ctx, stmt)
	if err != nil {
		logger.Error("failed to prepare named context", err)
		return errors.Wrap(err, op)
	}

	arg := map[string]interface{}{
		"uid":     uid,
		"card_id": cardID,
	}

	if err := namedStmt.QueryRowxContext(ctx, arg).Err(); err != nil {
		logger.Error("failed to query row context", err)
		return errors.Wrap(err, op)
	}

	return nil
}
