package handlers

import (
	registrierung2 "cqrs-playground/bibliothek/nutzer/registrierung"
	"encoding/json"
	"net/http"
)

type NutzerRegistrierungAPI struct {
	NutzerRegistrierungHandler *registrierung2.NutzerRegistrierungHandler
}

func NewNutzerRegistrierungAPI(
	nutzerRegistrierungHandler *registrierung2.NutzerRegistrierungHandler,
) *NutzerRegistrierungAPI {
	return &NutzerRegistrierungAPI{
		NutzerRegistrierungHandler: nutzerRegistrierungHandler,
	}
}

func (api *NutzerRegistrierungAPI) RegistriereNutzer(w http.ResponseWriter, r *http.Request) {
	var cmd registrierung2.NutzerRegistrierungCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := api.NutzerRegistrierungHandler.Handle(cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "neuer benutzer registriert"})
}
