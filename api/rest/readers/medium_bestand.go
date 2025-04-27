package readers

import (
	"context"
	"cqrs-playground/bibliothek/medien/projections/bestand"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type MediumBestandAPI struct {
	db            *pgxpool.Pool
	bestandReader *bestand.MediumBestandReader
	ctx           context.Context
	cancel        context.CancelFunc
}

func NewMediumBestandAPI(db *pgxpool.Pool) *MediumBestandAPI {
	ctx, cancel := context.WithCancel(context.Background())
	br := bestand.NewMediumBestandReader(db)
	return &MediumBestandAPI{
		db:            db,
		bestandReader: br,
		ctx:           ctx,
		cancel:        cancel,
	}
}

func (m *MediumBestandAPI) GetAll(w http.ResponseWriter, r *http.Request) {
	bestand, err := m.bestandReader.GetAll(m.ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bestand)
}
