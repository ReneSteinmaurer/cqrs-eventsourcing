package handlers

import (
	"cqrs-playground/api/rest"
	"cqrs-playground/bibliothek/medien/ausleihen"
	"cqrs-playground/bibliothek/medien/erwerben"
	"cqrs-playground/bibliothek/medien/katalogisieren"
	"cqrs-playground/bibliothek/medien/rueckgeben"
	"encoding/json"
	"net/http"
)

type MediumHandlerAPI struct {
	ErwerbeMediumHandler       *erwerben.ErwerbeMediumHandler
	KatalogisiereMediumHandler *katalogisieren.KatalogisiereMediumHandler
	VerleiheMediumHandler      *ausleihen.VerleiheMediumHandler
	RueckgebenMediumHandler    *rueckgeben.MediumRueckgabeHandler
}

func NewErwerbeMediumAPI(
	erwerbeMediumHandler *erwerben.ErwerbeMediumHandler,
	katalogisiereMediumHandler *katalogisieren.KatalogisiereMediumHandler,
	verleiheMediumHandler *ausleihen.VerleiheMediumHandler,
	rueckgebenMediumHandler *rueckgeben.MediumRueckgabeHandler,
) *MediumHandlerAPI {
	return &MediumHandlerAPI{
		ErwerbeMediumHandler:       erwerbeMediumHandler,
		KatalogisiereMediumHandler: katalogisiereMediumHandler,
		VerleiheMediumHandler:      verleiheMediumHandler,
		RueckgebenMediumHandler:    rueckgebenMediumHandler,
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
	json.NewEncoder(w).Encode(map[string]string{"status": "neues medium hinzugefügt"})
}

func (api *MediumHandlerAPI) KatalogisiereMedium(w http.ResponseWriter, r *http.Request) {
	var cmd katalogisieren.KatalogisiereMediumCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		rest.SendResponseErrors(&w, err.Error())
		return
	}

	aggregateId, err := api.KatalogisiereMediumHandler.Handle(cmd)
	if err != nil {
		rest.SendResponseErrors(&w, err.Error())
		return
	}

	res := rest.NewResponseContentMessage("medium katalogisiert", aggregateId)
	rest.SendResponse(res, &w)
}

func (api *MediumHandlerAPI) VerleiheMedium(w http.ResponseWriter, r *http.Request) {
	var cmd ausleihen.VerleiheMediumCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		rest.SendResponseErrors(&w, err.Error())
		return
	}

	aggregateId, err := api.VerleiheMediumHandler.Handle(cmd)
	if err != nil {
		rest.SendResponseErrors(&w, err.Error())
		return
	}

	res := rest.NewResponseContentMessage("medium verliehen", aggregateId)
	rest.SendResponse(res, &w)
}

func (api *MediumHandlerAPI) GebeMediumZurueck(w http.ResponseWriter, r *http.Request) {
	var cmd rueckgeben.MediumRueckgebenCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	aggregateId, err := api.RueckgebenMediumHandler.Handle(cmd)
	if err != nil {
		rest.SendResponseErrors(&w, err.Error())
		return
	}

	res := rest.NewResponseContentMessage("medium zurueckgegeben", aggregateId)
	rest.SendResponse(res, &w)
}
