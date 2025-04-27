package verliehen

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type MediumVerliehen struct {
	MediumId     string    `db:"medium_id"`
	VerliehenVon time.Time `db:"verliehen_von"`
	VerliehenBis time.Time `db:"verliehen_bis"`
	NutzerId     string    `db:"nutzer_id"`
}

type MediumVerliehenReader struct {
	db *pgxpool.Pool
}

func NewMediumVerliehenReader(db *pgxpool.Pool) *MediumVerliehenReader {
	return &MediumVerliehenReader{db: db}
}

func (r *MediumVerliehenReader) GetByMediumId(ctx context.Context, mediumId string) (*MediumVerliehen, error) {
	var medium MediumVerliehen
	err := r.db.QueryRow(ctx, `
		SELECT (medium_id, verliehen_von, verliehen_bis, nutzer_id)
		FROM medium_verliehen WHERE medium_id = $1
	`, mediumId).Scan(&medium)
	if err != nil {
		return nil, err
	}
	return &medium, nil
}
