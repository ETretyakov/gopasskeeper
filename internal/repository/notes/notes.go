package notes

import (
	"context"
	"gopasskeeper/internal/domain/models"
	"gopasskeeper/internal/logger"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// NotesRepoImpl is a structure to implement notes repository.
type NotesRepoImpl struct {
	db *sqlx.DB
}

// New is a builder function for NotesRepoImpl.
func New(db *sqlx.DB) *NotesRepoImpl {
	return &NotesRepoImpl{db: db}
}

// Add is a repository method to add record.
func (r *NotesRepoImpl) Add(
	ctx context.Context,
	uid string,
	name string,
	content string,
	meta string,
) (string, error) {
	const op = "repository.Notes.Add"

	stmt := `
	INSERT INTO public.sec_notes(uid, name, content, meta)
	VALUES (:uid, :name, :content, :meta)
	RETURNING id;
	`

	namedStmt, err := r.db.PrepareNamedContext(ctx, stmt)
	if err != nil {
		logger.Error("failed to prepare named context", err)
		return "", errors.Wrap(err, op)
	}

	arg := map[string]interface{}{
		"uid":     uid,
		"name":    name,
		"content": content,
		"meta":    meta,
	}

	var noteID string
	if err := namedStmt.QueryRowxContext(ctx, arg).Scan(&noteID); err != nil {
		logger.Error("failed to query row context", err)
		return "", errors.Wrap(err, op)
	}

	return noteID, nil
}

// GetSecret is a repository method to get record secret.
func (r *NotesRepoImpl) GetSecret(
	ctx context.Context,
	uid string,
	noteID string,
) (*models.NoteSecret, error) {
	const op = "repository.Notes.GetSecret"

	stmt := `
	SELECT sn.name    AS "name",
		   sn.content AS "content",
		   sn.meta    AS "meta"
	FROM sec_notes sn
	WHERE sn.uid = :uid AND 
	      sn.id  = :note_id
	LIMIT 1;
	`

	namedStmt, err := r.db.PrepareNamedContext(ctx, stmt)
	if err != nil {
		logger.Error("failed to prepare named context", err)
		return nil, errors.Wrap(err, op)
	}

	arg := map[string]interface{}{
		"uid":     uid,
		"note_id": noteID,
	}

	var noteSecret models.NoteSecret
	if err := namedStmt.QueryRowxContext(ctx, arg).StructScan(&noteSecret); err != nil {
		logger.Error("failed to query row context", err)
		return nil, errors.Wrap(err, op)
	}

	return &noteSecret, nil
}

// Search is a repository method to search records.
func (r *NotesRepoImpl) Search(
	ctx context.Context,
	uid string,
	schema *models.NoteSearchRequest,
) (*models.NoteSearchResponse, error) {
	const op = "repository.Notes.Search"

	// selecting items
	stmt := `
	SELECT sn.id   AS "id",
	       sn.name AS "name"
	FROM sec_notes sn
	WHERE sn.uid = :uid AND
		  sn.name ILIKE :substring
	ORDER BY sn.name
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

	var items []*models.NoteSearchItem
	for rows.Next() {
		var noteSearchItem models.NoteSearchItem
		if err := rows.StructScan(&noteSearchItem); err != nil {
			logger.Error("failed to query row context", err)
			return nil, errors.Wrap(err, op)
		}

		items = append(items, &noteSearchItem)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, op)
	}

	// selecting count
	stmt = `
	SELECT count(*) AS "count"
	FROM sec_notes sn
	WHERE sn.uid = :uid AND
		  sn.name ILIKE :substring;
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

	return &models.NoteSearchResponse{
		Count: count,
		Items: items,
	}, nil
}

// Remove is a repository method to remove a record.
func (r *NotesRepoImpl) Remove(
	ctx context.Context,
	uid string,
	noteID string,
) error {
	const op = "repository.Notes.Remove"

	stmt := `
	DELETE FROM sec_notes
	WHERE uid = :uid AND
	      id  = :note_id;
	`

	namedStmt, err := r.db.PrepareNamedContext(ctx, stmt)
	if err != nil {
		logger.Error("failed to prepare named context", err)
		return errors.Wrap(err, op)
	}

	arg := map[string]interface{}{
		"uid":     uid,
		"note_id": noteID,
	}

	if err := namedStmt.QueryRowxContext(ctx, arg).Err(); err != nil {
		logger.Error("failed to query row context", err)
		return errors.Wrap(err, op)
	}

	return nil
}
