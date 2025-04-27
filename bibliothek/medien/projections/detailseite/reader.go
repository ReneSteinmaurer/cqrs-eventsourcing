package detailseite

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type MediumDetail struct {
	MediumID         string
	ISBN             string
	Titel            string
	Genre            string
	Typ              string
	Standort         string
	Signatur         string
	ExemplarCode     string
	AktuellVerliehen bool
	VerliehenAn      *string
	VerliehenVon     *time.Time
	FaelligBis       *time.Time
	Status           string

	Historie []MediumHistorieEntry
}

type MediumHistorieEntry struct {
	EventType string
	Timestamp time.Time
	Payload   json.RawMessage
}

type DetailseiteReader struct {
	db *pgxpool.Pool
}

func NewDetailseiteReader(db *pgxpool.Pool) *DetailseiteReader {
	return &DetailseiteReader{db: db}
}

func (r *DetailseiteReader) GetMediumDetailWithHistorie(ctx context.Context, mediumID string) (*MediumDetail, error) {
	const detailQuery = `
		SELECT medium_id, isbn, titel, genre, typ, standort, signatur, exemplar_code, aktuell_verliehen, verliehen_an, verliehen_von, faellig_bis, status
		FROM medium_details
		WHERE medium_id = $1
	`
	var detail MediumDetail
	err := r.db.QueryRow(ctx, detailQuery, mediumID).Scan(
		&detail.MediumID,
		&detail.ISBN,
		&detail.Titel,
		&detail.Genre,
		&detail.Typ,
		&detail.Standort,
		&detail.Signatur,
		&detail.ExemplarCode,
		&detail.AktuellVerliehen,
		&detail.VerliehenAn,
		&detail.VerliehenVon,
		&detail.FaelligBis,
		&detail.Status,
	)
	if err != nil {
		return nil, err
	}

	const historyQuery = `
		SELECT event_type, event_timestamp, payload
		FROM medium_historie
		WHERE medium_id = $1
		ORDER BY event_timestamp ASC
	`
	rows, err := r.db.Query(ctx, historyQuery, mediumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry MediumHistorieEntry
		err := rows.Scan(&entry.EventType, &entry.Timestamp, &entry.Payload)
		if err != nil {
			return nil, err
		}
		detail.Historie = append(detail.Historie, entry)
	}

	return &detail, nil
}
