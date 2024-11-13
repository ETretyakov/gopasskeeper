package accounts

import (
	"context"
	"gopasskeeper/internal/domain/models"
	"gopasskeeper/internal/logger"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// AccountsRepoImpl is a structure to implement accounts repository.
type AccountsRepoImpl struct {
	db *sqlx.DB
}

// New is a builder function for AccountsRepoImpl.
func New(db *sqlx.DB) *AccountsRepoImpl {
	return &AccountsRepoImpl{db: db}
}

// Add is a repository method to add record.
func (r *AccountsRepoImpl) Add(
	ctx context.Context,
	uid string,
	login string,
	server string,
	password string,
) (string, error) {
	const op = "repository.Accounts.Add"

	stmt := `
	INSERT INTO public.sec_accounts(uid, login, password, server)
	VALUES (:uid, :login, :password, :server)
	RETURNING id;
	`

	namedStmt, err := r.db.PrepareNamedContext(ctx, stmt)
	if err != nil {
		logger.Error("failed to prepare named context", err)
		return "", errors.Wrap(err, op)
	}

	arg := map[string]interface{}{
		"uid":      uid,
		"login":    login,
		"password": password,
		"server":   server,
	}

	var accountID string
	if err := namedStmt.QueryRowxContext(ctx, arg).Scan(&accountID); err != nil {
		logger.Error("failed to query row context", err)
		return "", errors.Wrap(err, op)
	}

	return accountID, nil
}

// GetSecret is a repository method to get record secret.
func (r *AccountsRepoImpl) GetSecret(
	ctx context.Context,
	uid string,
	accountID string,
) (*models.AccountSecret, error) {
	const op = "repository.Accounts.GetSecret"

	stmt := `
	SELECT sa.login    AS "login",
		   sa.server   AS "server",
		   sa.password AS "password"
	FROM sec_accounts sa
	WHERE sa.uid = :uid AND 
	      sa.id  = :account_id
	LIMIT 1;
	`

	namedStmt, err := r.db.PrepareNamedContext(ctx, stmt)
	if err != nil {
		logger.Error("failed to prepare named context", err)
		return nil, errors.Wrap(err, op)
	}

	arg := map[string]interface{}{
		"uid":        uid,
		"account_id": accountID,
	}

	var accountSecret models.AccountSecret
	if err := namedStmt.QueryRowxContext(ctx, arg).StructScan(&accountSecret); err != nil {
		logger.Error("failed to query row context", err)
		return nil, errors.Wrap(err, op)
	}

	return &accountSecret, nil
}

// Search is a repository method to search records.
func (r *AccountsRepoImpl) Search(
	ctx context.Context,
	uid string,
	schema *models.AccountSearchRequest,
) (*models.AccountSearchResponse, error) {
	const op = "repository.Accounts.Search"

	// selecting items
	stmt := `
	SELECT sa.id       AS "id",
	       sa.login    AS "login",
		   sa.server   AS "server"
	FROM sec_accounts sa
	WHERE sa.uid = :uid AND
		  (sa.server ILIKE :substring OR
		   sa.login  ILIKE :substring)
	ORDER BY sa.server, sa.login
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

	var items []*models.AccountSearchItem
	for rows.Next() {
		var accountSearchItem models.AccountSearchItem
		if err := rows.StructScan(&accountSearchItem); err != nil {
			logger.Error("failed to query row context", err)
			return nil, errors.Wrap(err, op)
		}

		items = append(items, &accountSearchItem)
	}

	// selecting count
	stmt = `
	SELECT count(*) AS "count"
	FROM sec_accounts sa
	WHERE sa.uid = :uid AND
		  (sa.server ILIKE :substring OR
		   sa.login  ILIKE :substring);
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

	return &models.AccountSearchResponse{
		Count: count,
		Items: items,
	}, nil
}

// Remove is a repository method to remove a record.
func (r *AccountsRepoImpl) Remove(
	ctx context.Context,
	uid string,
	accountID string,
) error {
	const op = "repository.Accounts.Remove"

	stmt := `
	DELETE FROM sec_accounts
	WHERE uid = :uid AND
	      id  = :account_id;
	`

	namedStmt, err := r.db.PrepareNamedContext(ctx, stmt)
	if err != nil {
		logger.Error("failed to prepare named context", err)
		return errors.Wrap(err, op)
	}

	arg := map[string]interface{}{
		"uid":        uid,
		"account_id": accountID,
	}

	if err := namedStmt.QueryRowxContext(ctx, arg).Err(); err != nil {
		logger.Error("failed to query row context", err)
		return errors.Wrap(err, op)
	}

	return nil
}
