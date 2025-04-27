package detailseite

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type MediumHistorieEntry struct {
	EventType string          `json:"eventType"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}

type MediumDetail struct {
	MediumID          string     `json:"mediumId"`
	ISBN              string     `json:"isbn"`
	Titel             string     `json:"titel"`
	Genre             string     `json:"genre"`
	Typ               string     `json:"typ"`
	Standort          *string    `json:"standort"`
	Signatur          *string    `json:"signatur"`
	ExemplarCode      *string    `json:"exemplarCode"`
	AktuellVerliehen  bool       `json:"aktuellVerliehen"`
	VerliehenAn       *string    `json:"verliehenAn"`
	VerliehenVon      *time.Time `json:"verliehenVon"`
	FaelligBis        *time.Time `json:"faelligBis"`
	Status            string     `json:"status"`
	ErworbenAm        *time.Time `json:"erworbenAm"`
	KatalogisiertAm   *time.Time `json:"katalogisiertAm"`
	VerliehenNutzerId *string    `json:"verliehenNutzerId"`

	Historie []MediumHistorieEntry `json:"historie"`
}

type DetailseiteReader struct {
	db *pgxpool.Pool
}

func NewDetailseiteReader(db *pgxpool.Pool) *DetailseiteReader {
	return &DetailseiteReader{db: db}
}

func (r *DetailseiteReader) GetMediumDetailWithHistorie(ctx context.Context, mediumID string) (*MediumDetail, error) {
	const detailQuery = `
		SELECT medium_id, isbn, titel, genre, typ, standort, signatur, exemplar_code, aktuell_verliehen, verliehen_an, verliehen_von, faellig_bis, status, erworben_am, katalogisiert_am, verliehen_an_nutzer_id
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
		&detail.ErworbenAm,
		&detail.KatalogisiertAm,
		&detail.VerliehenNutzerId,
	)
	if err != nil {
		return nil, err
	}

	const historyQuery = `
		SELECT event_type, event_timestamp, payload
		FROM medium_historie
		WHERE medium_id = $1
		ORDER BY event_timestamp desc
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
