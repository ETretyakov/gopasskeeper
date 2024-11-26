package files

import (
	"context"
	"gopasskeeper/internal/domain/models"
	"gopasskeeper/internal/logger"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// FilesRepoImpl is a structure to implement files repository.
type FilesRepoImpl struct {
	db *sqlx.DB
}

// New is a builder function for FilesRepoImpl.
func New(db *sqlx.DB) *FilesRepoImpl {
	return &FilesRepoImpl{db: db}
}

// Add is a repository method to add record.
func (r *FilesRepoImpl) Add(
	ctx context.Context,
	uid string,
	name string,
	meta string,
) (string, error) {
	const op = "repository.Files.Add"

	stmt := `
	INSERT INTO public.sec_files(uid, name, meta)
	VALUES (:uid, :name, :meta)
	RETURNING id;
	`

	namedStmt, err := r.db.PrepareNamedContext(ctx, stmt)
	if err != nil {
		logger.Error("failed to prepare named context", err)
		return "", errors.Wrap(err, op)
	}

	arg := map[string]interface{}{
		"uid":  uid,
		"name": name,
		"meta": meta,
	}

	var fileID string
	if err := namedStmt.QueryRowxContext(ctx, arg).Scan(&fileID); err != nil {
		logger.Error("failed to query row context", err)
		return "", errors.Wrap(err, op)
	}

	return fileID, nil
}

// GetSecret is a repository method to get record secret.
func (r *FilesRepoImpl) GetSecret(
	ctx context.Context,
	uid string,
	fileID string,
) (string, string, error) {
	const op = "repository.Files.GetSecret"

	stmt := `
	SELECT sf.name   AS "name",
	       sf.meta   AS "meta"
	FROM sec_files sf
	WHERE sf.uid = :uid AND 
	      sf.id  = :file_id
	LIMIT 1;
	`

	namedStmt, err := r.db.PrepareNamedContext(ctx, stmt)
	if err != nil {
		logger.Error("failed to prepare named context", err)
		return "", "", errors.Wrap(err, op)
	}

	arg := map[string]interface{}{
		"uid":     uid,
		"file_id": fileID,
	}

	name := ""
	meta := ""
	if err := namedStmt.QueryRowxContext(ctx, arg).Scan(&name, &meta); err != nil {
		logger.Error("failed to query row context", err)
		return "", "", errors.Wrap(err, op)
	}

	return name, meta, nil
}

// Search is a repository method to search records.
func (r *FilesRepoImpl) Search(
	ctx context.Context,
	uid string,
	schema *models.FileSearchRequest,
) (*models.FileSearchResponse, error) {
	const op = "repository.Files.Search"

	// selecting items
	stmt := `
	SELECT sf.id   AS "id",
	       sf.name AS "name"
	FROM sec_files sf
	WHERE sf.uid = :uid AND
		  sf.name ILIKE :substring
	ORDER BY sf.name
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

	var items []*models.FileSearchItem
	for rows.Next() {
		var fileSearchItem models.FileSearchItem
		if err := rows.StructScan(&fileSearchItem); err != nil {
			logger.Error("failed to query row context", err)
			return nil, errors.Wrap(err, op)
		}

		items = append(items, &fileSearchItem)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, op)
	}

	// selecting count
	stmt = `
	SELECT count(*) AS "count"
	FROM sec_files sf
	WHERE sf.uid = :uid AND
		  sf.name ILIKE :substring;
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

	return &models.FileSearchResponse{
		Count: count,
		Items: items,
	}, nil
}

// Remove is a repository method to remove a record.
func (r *FilesRepoImpl) Remove(
	ctx context.Context,
	uid string,
	fileID string,
) error {
	const op = "repository.Files.Remove"

	stmt := `
	DELETE FROM sec_files
	WHERE uid = :uid AND
	      id  = :file_id;
	`

	namedStmt, err := r.db.PrepareNamedContext(ctx, stmt)
	if err != nil {
		logger.Error("failed to prepare named context", err)
		return errors.Wrap(err, op)
	}

	arg := map[string]interface{}{
		"uid":     uid,
		"file_id": fileID,
	}

	if err := namedStmt.QueryRowxContext(ctx, arg).Err(); err != nil {
		logger.Error("failed to query row context", err)
		return errors.Wrap(err, op)
	}

	return nil
}
