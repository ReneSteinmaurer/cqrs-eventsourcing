package nutzer

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Nutzer struct {
	NutzerId string `db:"nutzer_id" json:"nutzerId"`
	Email    string `db:"email" json:"email"`
	Vorname  string `db:"vorname" json:"vorname"`
	Nachname string `db:"nachname" json:"nachname"`
}

type NutzerReader struct {
	db *pgxpool.Pool
}

func NewNutzerReader(db *pgxpool.Pool) *NutzerReader {
	return &NutzerReader{db: db}
}

func (r *NutzerReader) Exists(ctx context.Context, nutzerId string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM nutzer WHERE nutzer_id = $1)", nutzerId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *NutzerReader) GetNutzer(ctx context.Context, nutzerId string) (Nutzer, error) {
	var nutzer Nutzer
	err := r.db.QueryRow(ctx, `
		SELECT 
			nutzer_id,
			email,
			vorname,
			nachname 
			FROM nutzer 
		WHERE nutzer_id = $1`, nutzerId).Scan(&nutzer.NutzerId, &nutzer.Email, &nutzer.Vorname, &nutzer.Nachname)
	return nutzer, err
}

func (r *NutzerReader) FindNutzerByEmailPrefix(ctx context.Context, prefix string) ([]Nutzer, error) {
	const query = `
		SELECT nutzer_id, email, vorname, nachname
		FROM nutzer
		WHERE email ILIKE $1
		ORDER BY email ASC
		LIMIT 20
	`

	rows, err := r.db.Query(ctx, query, prefix+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nutzerList := make([]Nutzer, 0)
	for rows.Next() {
		var nutzer Nutzer
		err := rows.Scan(&nutzer.NutzerId, &nutzer.Email, &nutzer.Vorname, &nutzer.Nachname)
		if err != nil {
			return nil, err
		}
		nutzerList = append(nutzerList, nutzer)
	}

	return nutzerList, nil
}
