package handlers

import (
	"cqrs-playground/api/rest"
	"cqrs-playground/bibliothek/medien/ausleihen"
	"cqrs-playground/bibliothek/medien/erwerben"
	"cqrs-playground/bibliothek/medien/katalogisieren"
	"cqrs-playground/bibliothek/medien/rueckgeben"
	"cqrs-playground/bibliothek/medien/verlieren/bestands_verlust"
	"cqrs-playground/bibliothek/medien/verlieren/verloren_duch_benutzer"
	bestands_verlust2 "cqrs-playground/bibliothek/medien/wiederaufgefunden/bestands_verlust"
	"encoding/json"
	"net/http"
)

type MediumHandlerAPI struct {
	ErwerbeMediumHandler           *erwerben.ErwerbeMediumHandler
	KatalogisiereMediumHandler     *katalogisieren.KatalogisiereMediumHandler
	VerleiheMediumHandler          *ausleihen.VerleiheMediumHandler
	RueckgebenMediumHandler        *rueckgeben.MediumRueckgabeHandler
	VerlorenDurchNutzerHandler     *verloren_duch_benutzer.MediumVerlorenDurchBenutzerHandler
	BestandsverlustHandler         *bestands_verlust.MediumBestandsverlustHandler
	BestandsverlustAufhebenHandler *bestands_verlust2.MediumBestandsverlustAufhebenHandler
}

func NewErwerbeMediumAPI(
	erwerbeMediumHandler *erwerben.ErwerbeMediumHandler,
	katalogisiereMediumHandler *katalogisieren.KatalogisiereMediumHandler,
	verleiheMediumHandler *ausleihen.VerleiheMediumHandler,
	rueckgebenMediumHandler *rueckgeben.MediumRueckgabeHandler,
	verlorenDurchNutzerHandler *verloren_duch_benutzer.MediumVerlorenDurchBenutzerHandler,
	bestandsverlustHandler *bestands_verlust.MediumBestandsverlustHandler,
	bestandsverlustAufhebenHandler *bestands_verlust2.MediumBestandsverlustAufhebenHandler,
) *MediumHandlerAPI {
	return &MediumHandlerAPI{
		ErwerbeMediumHandler:           erwerbeMediumHandler,
		KatalogisiereMediumHandler:     katalogisiereMediumHandler,
		VerleiheMediumHandler:          verleiheMediumHandler,
		RueckgebenMediumHandler:        rueckgebenMediumHandler,
		VerlorenDurchNutzerHandler:     verlorenDurchNutzerHandler,
		BestandsverlustHandler:         bestandsverlustHandler,
		BestandsverlustAufhebenHandler: bestandsverlustAufhebenHandler,
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

func (api *MediumHandlerAPI) MediumVerlorenVonNutzer(w http.ResponseWriter, r *http.Request) {
	var cmd verloren_duch_benutzer.MediumVerlorenDurchBenutzerCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	aggregateId, err := api.VerlorenDurchNutzerHandler.Handle(cmd)
	if err != nil {
		rest.SendResponseErrors(&w, err.Error())
		return
	}

	res := rest.NewResponseContentMessage("medium verloren", aggregateId)
	rest.SendResponse(res, &w)
}

func (api *MediumHandlerAPI) MediumBestandsverlust(w http.ResponseWriter, r *http.Request) {
	var cmd bestands_verlust.MediumBestandsverlustCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	aggregateId, err := api.BestandsverlustHandler.Handle(cmd)
	if err != nil {
		rest.SendResponseErrors(&w, err.Error())
		return
	}

	res := rest.NewResponseContentMessage("medium verloren", aggregateId)
	rest.SendResponse(res, &w)
}

func (api *MediumHandlerAPI) MediumBestandsverlustAufheben(w http.ResponseWriter, r *http.Request) {
	var cmd bestands_verlust2.MediumBestandsverlustAufhebenCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	aggregateId, err := api.BestandsverlustAufhebenHandler.Handle(cmd)
	if err != nil {
		rest.SendResponseErrors(&w, err.Error())
		return
	}

	res := rest.NewResponseContentMessage("medium verloren", aggregateId)
	rest.SendResponse(res, &w)
}
