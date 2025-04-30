package detailseite

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type DetailseiteReader struct {
	db *pgxpool.Pool
}

func NewDetailseiteReader(db *pgxpool.Pool) *DetailseiteReader {
	return &DetailseiteReader{db: db}
}

func (r *DetailseiteReader) GetNutzerDetails(ctx context.Context, nutzerId string) (*NutzerDetails, error) {
	const query = `
        SELECT vorname, nachname, status, registriert_am, aktive_ausleihen, letzte_notizen, verlorene_medien, sperrgrund
        FROM nutzer_details
        WHERE nutzer_id = $1
    `

	var (
		vorname, nachname, status                       string
		registriertAm                                   time.Time
		ausleihenJSON, notizenJSON, verloreneMedienJSON []byte
		sperrgrund                                      *string
	)

	err := r.db.QueryRow(ctx, query, nutzerId).Scan(
		&vorname,
		&nachname,
		&status,
		&registriertAm,
		&ausleihenJSON,
		&notizenJSON,
		&verloreneMedienJSON,
		&sperrgrund,
	)
	if err != nil {
		return nil, err
	}

	var ausleihen []AktiveAusleihe
	if err := json.Unmarshal(ausleihenJSON, &ausleihen); err != nil {
		return nil, err
	}

	var notizen []Notiz
	if err := json.Unmarshal(notizenJSON, &notizen); err != nil {
		return nil, err
	}

	var verloreneMedien []VerlorenesMedium
	if err := json.Unmarshal(verloreneMedienJSON, &verloreneMedien); err != nil {
		return nil, err
	}

	return &NutzerDetails{
		NutzerId:        nutzerId,
		Vorname:         vorname,
		Nachname:        nachname,
		Status:          NutzerStatus(status),
		RegistriertAm:   &registriertAm,
		AktiveAusleihen: ausleihen,
		LetzteNotizen:   notizen,
		VerloreneMedien: verloreneMedien,
		Sperrgrund:      sperrgrund,
	}, nil
}
