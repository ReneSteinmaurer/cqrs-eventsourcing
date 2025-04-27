package readers

import (
	"context"
	"cqrs-playground/bibliothek/medien/projections/detailseite"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type MediumDetailsAPI struct {
	db                *pgxpool.Pool
	detailseiteReader *detailseite.DetailseiteReader
	ctx               context.Context
	cancel            context.CancelFunc
}

func NewMediumDetailsAPI(db *pgxpool.Pool) *MediumDetailsAPI {
	ctx, cancel := context.WithCancel(context.Background())
	br := detailseite.NewDetailseiteReader(db)
	return &MediumDetailsAPI{
		detailseiteReader: br,
		db:                db,
		ctx:               ctx,
		cancel:            cancel,
	}
}

func (m *MediumDetailsAPI) GetAll(w http.ResponseWriter, r *http.Request) {
	mediumID := r.URL.Query().Get("mediumId")
	mediumDetails, err := m.detailseiteReader.GetMediumDetailWithHistorie(m.ctx, mediumID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mediumDetails)
}
