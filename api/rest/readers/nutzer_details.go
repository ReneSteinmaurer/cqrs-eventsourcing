package readers

import (
	"context"
	"cqrs-playground/bibliothek/nutzer/projections/detailseite"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type NutzerDetailsAPI struct {
	db                *pgxpool.Pool
	detailseiteReader *detailseite.DetailseiteReader
	ctx               context.Context
	cancel            context.CancelFunc
}

func NewNutzerDetailsAPI(db *pgxpool.Pool) *NutzerDetailsAPI {
	ctx, cancel := context.WithCancel(context.Background())
	br := detailseite.NewDetailseiteReader(db)
	return &NutzerDetailsAPI{
		detailseiteReader: br,
		db:                db,
		ctx:               ctx,
		cancel:            cancel,
	}
}

func (m *NutzerDetailsAPI) GetAll(w http.ResponseWriter, r *http.Request) {
	mediumID := r.URL.Query().Get("nutzerId")
	mediumDetails, err := m.detailseiteReader.GetNutzerDetailWithHistorie(m.ctx, mediumID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mediumDetails)
}
