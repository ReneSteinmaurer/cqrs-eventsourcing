package isbn_index

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ISBNIndexReader struct {
	db *pgxpool.Pool
}

func NewPostgresISBNIndexReader(db *pgxpool.Pool) *ISBNIndexReader {
	return &ISBNIndexReader{db: db}
}

func (r *ISBNIndexReader) Exists(ctx context.Context, isbn string) (bool, error) {
	var id string
	err := r.db.QueryRow(ctx, "SELECT medium_id FROM isbn_index WHERE isbn = $1", isbn).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *ISBNIndexReader) GetMediumId(ctx context.Context, isbn string) (string, error) {
	var id string
	err := r.db.QueryRow(ctx, "SELECT medium_id FROM isbn_index WHERE isbn = $1", isbn).Scan(&id)
	return id, err
}

func (r *ISBNIndexReader) GetISBN(ctx context.Context, mediumId string) (string, error) {
	var id string
	err := r.db.QueryRow(ctx, "SELECT isbn FROM isbn_index WHERE medium_id = $1", mediumId).Scan(&id)
	return id, err
}
