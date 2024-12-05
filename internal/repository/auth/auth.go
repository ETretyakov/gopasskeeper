package auth

import (
	"context"
	"gopasskeeper/internal/domain/models"
	"gopasskeeper/internal/logger"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// AuthRepoImpl is an implementation for auth repository.
type AuthRepoImpl struct {
	db *sqlx.DB
}

// New is a builder function fot AuthRepoImpl.
func New(db *sqlx.DB) *AuthRepoImpl {
	return &AuthRepoImpl{db: db}
}

// SaveUser is a AuthRepoImpl method to save user.
func (r *AuthRepoImpl) SaveUser(
	ctx context.Context,
	login string,
	passHash []byte,
) (string, error) {
	const op = "repository.Auth.SaveUser"

	stmt := `
	INSERT INTO public.usr_users(login, pass_hash)
	VALUES (:login, :pass_hash)
	RETURNING id;
	`

	namedStmt, err := r.db.PrepareNamedContext(ctx, stmt)
	if err != nil {
		logger.Error("failed to prepare named context", err)
		return "", errors.Wrap(err, op)
	}

	arg := map[string]interface{}{
		"login":     login,
		"pass_hash": passHash,
	}

	var userID string
	if err := namedStmt.QueryRowxContext(ctx, arg).Scan(&userID); err != nil {
		logger.Error("failed to query row context", err)
		return "", errors.Wrap(err, op)
	}

	return userID, nil
}

// User is a AuthRepoImpl method to retrieve a user by login.
func (r *AuthRepoImpl) User(
	ctx context.Context,
	login string,
) (*models.UserAuth, error) {
	const op = "repository.Auth.User"

	stmt := `
	SELECT id, login, pass_hash
	FROM usr_users
	WHERE login = :login
	LIMIT 1;
	`

	namedStmt, err := r.db.PrepareNamedContext(ctx, stmt)
	if err != nil {
		logger.Error("failed to prepare named context", err)
		return nil, errors.Wrap(err, op)
	}

	arg := map[string]interface{}{
		"login": login,
	}

	user := models.UserAuth{}
	if err := namedStmt.QueryRowxContext(ctx, arg).StructScan(&user); err != nil {
		logger.Error("failed to query row context", err)
		return nil, errors.Wrap(err, op)
	}

	return &user, nil
}
