package rest

import (
	"cqrs-playground/bibliothek/medien/ausleihen"
	"cqrs-playground/bibliothek/medien/erwerben"
	"cqrs-playground/bibliothek/medien/katalogisieren"
	"encoding/json"
	"net/http"
)

type MediumHandlerAPI struct {
	ErwerbeMediumHandler       *erwerben.ErwerbeMediumHandler
	KatalogisiereMediumHandler *katalogisieren.KatalogisiereMediumHandler
	VerleiheMediumHandler      *ausleihen.VerleiheMediumHandler
}

func NewErwerbeMediumAPI(
	erwerbeMediumHandler *erwerben.ErwerbeMediumHandler,
	katalogisiereMediumHandler *katalogisieren.KatalogisiereMediumHandler,
	verleiheMediumHandler *ausleihen.VerleiheMediumHandler,
) *MediumHandlerAPI {
	return &MediumHandlerAPI{
		ErwerbeMediumHandler:       erwerbeMediumHandler,
		KatalogisiereMediumHandler: katalogisiereMediumHandler,
		VerleiheMediumHandler:      verleiheMediumHandler,
	}
}

func (api *MediumHandlerAPI) ErwerbeMedium(w http.ResponseWriter, r *http.Request) {
	var cmd erwerben.ErwerbeMediumCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := api.ErwerbeMediumHandler.Handle(cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "neues medium hinzugefuegt"})
}

func (api *MediumHandlerAPI) KatalogisiereMedium(w http.ResponseWriter, r *http.Request) {
	var cmd katalogisieren.KatalogisiereMediumCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := api.KatalogisiereMediumHandler.Handle(cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "medium katalogisiert"})
}

func (api *MediumHandlerAPI) VerleiheMedium(w http.ResponseWriter, r *http.Request) {
	var cmd ausleihen.VerleiheMediumCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := api.VerleiheMediumHandler.Handle(cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "medium verliehen"})
}
