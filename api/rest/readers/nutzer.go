package readers

import (
	"context"
	"cqrs-playground/api/rest"
	"cqrs-playground/bibliothek/nutzer/projections/nutzer"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type NutzerAPI struct {
	db           *pgxpool.Pool
	nutzerReader *nutzer.NutzerReader
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewNutzerAPI(db *pgxpool.Pool) *NutzerAPI {
	ctx, cancel := context.WithCancel(context.Background())
	nr := nutzer.NewNutzerReader(db)
	return &NutzerAPI{
		nutzerReader: nr,
		db:           db,
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (n *NutzerAPI) FindNutzerByEmailPrefix(w http.ResponseWriter, r *http.Request) {
	emailPrefix := r.URL.Query().Get("email")
	possibilities, err := n.nutzerReader.FindNutzerByEmailPrefix(n.ctx, emailPrefix)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := rest.NewResponseContentMessage("", possibilities)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
