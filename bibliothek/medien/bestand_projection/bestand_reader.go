package bestand_projection

import (
	"context"
	"cqrs-playground/bibliothek/medien/shared"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MediumBestand struct {
	MediumId     string            `db:"medium_id"`
	ISBN         string            `db:"isbn"`
	MediumType   shared.MediumType `db:"medium_type"`
	Name         string            `db:"name"`
	Genre        string            `db:"genre"`
	Signature    string            `db:"signature"`
	Standort     string            `db:"standort"`
	ExemplarCode string            `db:"exemplar_code"`
}

type MediumBestandReader struct {
	db *pgxpool.Pool
}

func NewMediumBestandReader(db *pgxpool.Pool) *MediumBestandReader {
	return &MediumBestandReader{db: db}
}

func (r *MediumBestandReader) GetByMediumId(ctx context.Context, mediumId string) (*MediumBestand, error) {
	var medium MediumBestand
	err := r.db.QueryRow(ctx, `
		SELECT (medium_id, isbn, medium_type, name, genre, signature, standort, exemplar_code)
		FROM medium_bestand WHERE medium_id = $1
	`, mediumId).Scan(&medium)
	if err != nil {
		return nil, err
	}
	return &medium, nil
}
