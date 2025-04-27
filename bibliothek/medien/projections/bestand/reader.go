package bestand

import (
	"context"
	"cqrs-playground/bibliothek/medien/erwerben/events"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MediumBestand struct {
	MediumId      string            `db:"medium_id" json:"mediumId"`
	ISBN          string            `db:"isbn" json:"ISBN"`
	MediumType    events.MediumType `db:"medium_type" json:"mediumType"`
	Name          string            `db:"name" json:"name"`
	Genre         string            `db:"genre" json:"genre"`
	Katalogisiert bool              `db:"katalogisiert" json:"katalogisiert"`
	Signature     string            `db:"signature" json:"signature"`
	Standort      string            `db:"standort" json:"standort"`
	ExemplarCode  string            `db:"exemplar_code" json:"exemplarCode"`
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
		SELECT medium_id, isbn, medium_type, name, genre, signature, standort, exemplar_code
		FROM medium_bestand
		WHERE medium_id = $1
	`, mediumId).Scan(
		&medium.MediumId,
		&medium.ISBN,
		&medium.MediumType,
		&medium.Name,
		&medium.Genre,
		&medium.Signature,
		&medium.Standort,
		&medium.ExemplarCode,
	)
	if err != nil {
		return nil, err
	}
	return &medium, nil
}

func (r *MediumBestandReader) GetAll(ctx context.Context) ([]MediumBestand, error) {
	rows, err := r.db.Query(ctx, `
		SELECT medium_id, 
			isbn,
			medium_type,
			name,
			genre,
			katalogisiert,
			coalesce(signature, '')    AS signature,
			coalesce(standort, '')     AS standort,
			coalesce(exemplar_code, '') AS exemplar_code
		FROM medium_bestand
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []MediumBestand
	for rows.Next() {
		var medium MediumBestand
		err := rows.Scan(
			&medium.MediumId,
			&medium.ISBN,
			&medium.MediumType,
			&medium.Name,
			&medium.Genre,
			&medium.Katalogisiert,
			&medium.Signature,
			&medium.Standort,
			&medium.ExemplarCode,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, medium)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}
