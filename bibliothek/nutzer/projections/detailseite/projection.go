package detailseite

import (
	"context"
	events4 "cqrs-playground/bibliothek/medien/ausleihen/events"
	events2 "cqrs-playground/bibliothek/medien/erwerben/events"
	"cqrs-playground/bibliothek/medien/projections/detailseite"
	"cqrs-playground/bibliothek/medien/projections/verliehen"
	events5 "cqrs-playground/bibliothek/medien/rueckgeben/events"
	"cqrs-playground/bibliothek/medien/verlieren/events"
	events6 "cqrs-playground/bibliothek/medien/wiederaufgefunden/events"
	events7 "cqrs-playground/bibliothek/nutzer/registrierung/events"
	"cqrs-playground/shared"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type DetailseiteProjection struct {
	EventStore          *shared.EventStore
	DB                  *pgxpool.Pool
	KafkaService        *shared.KafkaService
	notificationService *shared.NotificationService
	ctx                 context.Context
	cancel              context.CancelFunc
	mediumReader        *detailseite.DetailseiteReader
}

func NewDetailseiteProjection(
	ctx context.Context, eventStore *shared.EventStore, db *pgxpool.Pool, kafkaService *shared.KafkaService,
	notificationService *shared.NotificationService,
) *DetailseiteProjection {
	ctx, cancel := context.WithCancel(ctx)
	mr := detailseite.NewDetailseiteReader(db)
	return &DetailseiteProjection{
		ctx:                 ctx,
		cancel:              cancel,
		EventStore:          eventStore,
		notificationService: notificationService,
		DB:                  db,
		KafkaService:        kafkaService,
		mediumReader:        mr,
	}
}

func (d *DetailseiteProjection) Start() {
	go shared.ListenToEvent(d.ctx, d.KafkaService, events7.NutzerRegistriertEventType, d.applyNutzerRegistriert)
	go shared.ListenToEvent(d.ctx, d.KafkaService, events4.MediumVerliehenEventType, d.applyMediumVerliehen)
	go shared.ListenToEvent(d.ctx, d.KafkaService, events5.MediumZurueckgegebenEventType, d.applyMediumZurueckgegeben)
	go shared.ListenToEvent(d.ctx, d.KafkaService, events.MediumVerlorenDurchBenutzerEventType, d.applyMediumVerlorenDurchNutzer)
	go shared.ListenToEvent(d.ctx, d.KafkaService, events6.MediumWiederaufgefundenDurchNutzerEventType, d.applyMediumWiederaufgefundenDurchNutzer)
}

func (d *DetailseiteProjection) applyNutzerRegistriert(payloadJSON []byte) {
	var payload events7.NutzerRegistriertEvent
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		log.Println("Error unmarshalling event:", err)
		return
	}

	const query = `
		insert into nutzer_details (nutzer_id, vorname, nachname, status, registriert_am)
        values ($1, $2, $3, $4, now())
		on conflict do nothing 
	`

	_, err := d.DB.Exec(d.ctx, query, payload.NutzerId, payload.Vorname, payload.Nachname, NutzerstatusAktiv)
	if err != nil {
		log.Println("Error updating read-model:", err)
		return
	}

	if err := d.saveHistoryEvent(payload.NutzerId, events2.MediumErworbenEventType, payload); err != nil {
		log.Println("Error saving history event:", err)
		return
	}
	d.notificationService.NotifyProjectionUpdated(payload.NutzerId)
}

func (d *DetailseiteProjection) applyMediumVerliehen(payloadJSON []byte) {
	var payload verliehen.MediumVerliehen
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		log.Println("Error unmarshalling event:", err)
		return
	}

	medium, err := d.mediumReader.GetMediumDetails(d.ctx, payload.MediumId)
	if err != nil {
		log.Println("Error reading from mediumDetails table:", err)
		return
	}

	ausleihe := AktiveAusleihe{
		MediumId:      payload.MediumId,
		Titel:         medium.Titel,
		AusgeliehenAm: &payload.VerliehenVon,
		FaelligAm:     &payload.VerliehenBis,
		Status:        AusleihstatusAktiv,
	}

	ausleiheJSON, err := json.Marshal(ausleihe)
	if err != nil {
		log.Println("Error marshalling ausleihe:", err)
		return
	}

	const query = `
		UPDATE nutzer_details
		SET aktive_ausleihen = aktive_ausleihen || $1::jsonb
		WHERE nutzer_id = $2
	`
	_, err = d.DB.Exec(d.ctx, query, ausleiheJSON, payload.NutzerId)
	if err != nil {
		log.Println("Error updating aktive_ausleihen:", err)
		return
	}

	if err := d.saveHistoryEvent(payload.NutzerId, events2.MediumErworbenEventType, payload); err != nil {
		log.Println("Error saving history event:", err)
		return
	}
	d.notificationService.NotifyProjectionUpdated(payload.NutzerId)
}

func (d *DetailseiteProjection) applyMediumZurueckgegeben(payloadJSON []byte) {
	var payload events5.MediumZurueckgegebenEvent
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		log.Println("Error unmarshalling MediumZurueckgegeben:", err)
		return
	}

	const query = `
		UPDATE nutzer_details
		SET aktive_ausleihen = COALESCE((
		  SELECT jsonb_agg(elem)
		  FROM jsonb_array_elements(aktive_ausleihen) AS elem
		  WHERE elem->>'mediumId' <> $1
		), '[]'::jsonb)
		WHERE nutzer_id = $2;
	`
	_, err := d.DB.Exec(d.ctx, query, payload.MediumId, payload.NutzerId)
	if err != nil {
		log.Println("Error updating aktive_ausleihen on Rückgabe:", err)
		return
	}

	if err := d.saveHistoryEvent(payload.NutzerId, events5.MediumZurueckgegebenEventType, payload); err != nil {
		log.Println("Error saving history event (Rückgabe):", err)
		return
	}
	d.notificationService.NotifyProjectionUpdated(payload.NutzerId)
}

func (d *DetailseiteProjection) applyMediumVerlorenDurchNutzer(payloadJSON []byte) {
	var payload events.MediumVerlorenDurchBenutzerEvent
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		log.Println("Error unmarshalling MediumVerlorenDurchBenutzer:", err)
		return
	}

	medium, err := d.mediumReader.GetMediumDetails(d.ctx, payload.MediumId)
	if err != nil {
		log.Println("Error reading from mediumDetails table:", err)
		return
	}

	verloren := VerlorenesMedium{
		MediumId:      payload.MediumId,
		Titel:         medium.Titel,
		AusgeliehenAm: medium.VerliehenVon,
		FaelligAm:     medium.FaelligBis,
	}
	verlorenJSON, err := json.Marshal(verloren)
	if err != nil {
		log.Println("Error marshalling: verlorenes medium:", err)
		return
	}

	const query = `
		UPDATE nutzer_details
		SET aktive_ausleihen = COALESCE((
			SELECT jsonb_agg(elem)
			FROM jsonb_array_elements(aktive_ausleihen) AS elem
			WHERE elem->>'mediumId' <> $1
		), '[]'::jsonb),
		verlorene_medien = verlorene_medien || $2::jsonb
		WHERE nutzer_id = $3
	`

	_, err = d.DB.Exec(d.ctx, query, payload.MediumId, verlorenJSON, payload.NutzerId)
	if err != nil {
		log.Println("Error updating aktive_ausleihen on Verlust:", err)
		return
	}

	if err := d.saveHistoryEvent(payload.NutzerId, events.MediumVerlorenDurchBenutzerEventType, payload); err != nil {
		log.Println("Error saving history event (Verlust):", err)
		return
	}
	d.notificationService.NotifyProjectionUpdated(payload.NutzerId)
}

func (d *DetailseiteProjection) applyMediumWiederaufgefundenDurchNutzer(payloadJSON []byte) {
	var payload events6.MediumWiederaufgefundenDurchNutzerEvent
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		log.Println("Error unmarshalling MediumWiederaufgefundenDurchNutzer:", err)
		return
	}

	const query = `
		UPDATE nutzer_details
		SET verlorene_medien = COALESCE((
			SELECT jsonb_agg(elem)
			FROM jsonb_array_elements(verlorene_medien) AS elem
			WHERE elem->>'mediumId' <> $1
		), '[]'::jsonb)
		WHERE nutzer_id = $2
	`

	_, err := d.DB.Exec(d.ctx, query, payload.MediumId, payload.NutzerId)
	if err != nil {
		log.Println("Error updating verlorene_medien on Wiederaufgefunden:", err)
		return
	}

	if err := d.saveHistoryEvent(payload.NutzerId, events6.MediumWiederaufgefundenDurchNutzerEventType, payload); err != nil {
		log.Println("Error saving history event (Wiederaufgefunden):", err)
		return
	}

	d.notificationService.NotifyProjectionUpdated(payload.NutzerId)
}

func (d *DetailseiteProjection) saveHistoryEvent(nutzerID string, eventType string, eventPayload any) error {
	payloadJSON, err := json.Marshal(eventPayload)
	if err != nil {
		return err
	}

	const query = `
		INSERT INTO nutzer_history (nutzer_id, event_type, payload)
		VALUES ($1, $2, $3)
	`

	_, err = d.DB.Exec(d.ctx, query, nutzerID, eventType, payloadJSON)
	return err
}
