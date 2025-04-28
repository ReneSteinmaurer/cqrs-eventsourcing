package detailseite

import (
	"context"
	events4 "cqrs-playground/bibliothek/medien/ausleihen/events"
	events2 "cqrs-playground/bibliothek/medien/erwerben/events"
	events3 "cqrs-playground/bibliothek/medien/katalogisieren/events"
	events5 "cqrs-playground/bibliothek/medien/rueckgeben/events"
	"cqrs-playground/bibliothek/medien/verlieren/events"
	"cqrs-playground/bibliothek/nutzer/projections/nutzer"
	"cqrs-playground/shared"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

const (
	MediumErworben      = "ERWORBEN"
	MediumKatalogisiert = "KATALOGISIERT"
	MediumVerloren      = "VERLOREN"
	MediumVerliehen     = "VERLIEHEN"
)

type DetailseiteProjection struct {
	EventStore          *shared.EventStore
	DB                  *pgxpool.Pool
	KafkaService        *shared.KafkaService
	notificationService *shared.NotificationService
	ctx                 context.Context
	cancel              context.CancelFunc
	nutzerReader        *nutzer.NutzerReader
}

func NewDetailseiteProjection(
	ctx context.Context, eventStore *shared.EventStore, db *pgxpool.Pool, kafkaService *shared.KafkaService,
	notificationService *shared.NotificationService,
) *DetailseiteProjection {
	ctx, cancel := context.WithCancel(ctx)
	nutzerReader := nutzer.NewNutzerReader(db)
	return &DetailseiteProjection{
		ctx:                 ctx,
		cancel:              cancel,
		EventStore:          eventStore,
		notificationService: notificationService,
		DB:                  db,
		KafkaService:        kafkaService,
		nutzerReader:        nutzerReader,
	}
}

func (d *DetailseiteProjection) Start() {
	go shared.ListenToEvent(d.ctx, d.KafkaService, events2.MediumErworbenEventType, d.applyMediumErworben)
	go shared.ListenToEvent(d.ctx, d.KafkaService, events3.MediumKatalogisiertEventType, d.applyMediumKatalogisiert)
	go shared.ListenToEvent(d.ctx, d.KafkaService, events4.MediumVerliehenEventType, d.applyMediumVerliehen)
	go shared.ListenToEvent(d.ctx, d.KafkaService, events5.MediumZurueckgegebenEventType, d.applyMediumZurueckgegeben)
	go shared.ListenToEvent(d.ctx, d.KafkaService, events.MediumVerlorenDurchBenutzerEventType, d.applyMediumVerlorenDurchNutzer)
}

func (d *DetailseiteProjection) applyMediumErworben(payloadJSON []byte) {
	var payload events2.MediumErworbenEvent
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		log.Println("Error unmarshalling event:", err)
		return
	}

	const query = `
		insert into medium_details (medium_id, isbn, titel, genre, typ, status, erworben_am)
        values ($1, $2, $3, $4, $5, $6, now())
		on conflict do nothing 
	`

	_, err := d.DB.Exec(d.ctx, query, payload.MediumId, payload.ISBN, payload.Name, payload.Genre, payload.MediumType, MediumErworben)
	if err != nil {
		log.Println("Error updating read-model:", err)
		return
	}

	if err := d.saveHistoryEvent(payload.MediumId, events2.MediumErworbenEventType, payload); err != nil {
		log.Println("Error saving history event:", err)
		return
	}
	d.notificationService.Notify(payload.MediumId)
}

func (d *DetailseiteProjection) applyMediumKatalogisiert(payloadJSON []byte) {
	var payload events3.MediumKatalogisiertEvent
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		log.Println("Error unmarshalling event:", err)
		return
	}
	const query = `
		UPDATE medium_details
		SET 
			signatur = $1,
			standort = $2,
			exemplar_code = $3,
			status = $4,
			katalogisiert_am = now()
		WHERE medium_id = $5
	`

	_, err := d.DB.Exec(d.ctx, query, payload.Signature, payload.Standort, payload.ExemplarCode, MediumKatalogisiert, payload.MediumId)
	if err != nil {
		log.Println("Error updating read-model:", err)
		return
	}

	if err := d.saveHistoryEvent(payload.MediumId, events3.MediumKatalogisiertEventType, payload); err != nil {
		log.Println("Error saving history event:", err)
		return
	}
	d.notificationService.Notify(payload.MediumId)
}

func (d *DetailseiteProjection) applyMediumVerliehen(payloadJSON []byte) {
	var payload events4.MediumVerliehenEvent
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		log.Println("Error unmarshalling event:", err)
		return
	}

	nutzerModel, err := d.nutzerReader.GetNutzer(d.ctx, payload.NutzerId)
	if err != nil {
		log.Println("Error reading from nutzerModel table:", err)
		return
	}

	const query = `
		UPDATE medium_details
		SET 
		    verliehen_an = $1,
			verliehen_von = $2,
			faellig_bis = $3,
			aktuell_verliehen = true,
			status = $4,
			verliehen_an_nutzer_id = $5
		WHERE medium_id = $6
	`

	name := nutzerModel.Vorname + " " + nutzerModel.Nachname
	_, err = d.DB.Exec(d.ctx, query, name, payload.Von, payload.Bis, MediumVerliehen, payload.NutzerId, payload.MediumId)
	if err != nil {
		log.Println("Error updating read-model:", err)
		return
	}

	if err := d.saveHistoryEvent(payload.MediumId, events4.MediumVerliehenEventType, payload); err != nil {
		log.Println("Error saving history event:", err)
		return
	}
	d.notificationService.Notify(payload.MediumId)
}

func (d *DetailseiteProjection) applyMediumZurueckgegeben(payloadJSON []byte) {
	var payload events5.MediumZurueckgegebenEvent
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		log.Println("Error unmarshalling event:", err)
		return
	}

	const query = `
		UPDATE medium_details
		SET 
		    verliehen_an = $1,
			verliehen_von = $2,
			faellig_bis = $3,
			aktuell_verliehen = false,
			status = $4,
			verliehen_an_nutzer_id = $5
		WHERE medium_id = $6
	`

	_, err := d.DB.Exec(d.ctx, query, nil, nil, nil, MediumKatalogisiert, nil, payload.MediumId)
	if err != nil {
		log.Println("Error updating read-model:", err)
		return
	}

	if err := d.saveHistoryEvent(payload.MediumId, events5.MediumZurueckgegebenEventType, payload); err != nil {
		log.Println("Error saving history event:", err)
		return
	}
	d.notificationService.Notify(payload.MediumId)
}

func (d *DetailseiteProjection) applyMediumVerlorenDurchNutzer(payloadJSON []byte) {
	var payload events.MediumVerlorenDurchBenutzerEvent
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		log.Println("Error unmarshalling event:", err)
		return
	}

	nutzerModel, err := d.nutzerReader.GetNutzer(d.ctx, payload.NutzerId)
	if err != nil {
		log.Println("Error reading from nutzerModel table:", err)
		return
	}

	const query = `
		UPDATE medium_details
		SET 
		    verliehen_an = $1,
			verliehen_von = $2,
			faellig_bis = $3,
			aktuell_verliehen = false,
			status = $4,
			verliehen_an_nutzer_id = $5,
			verloren_am = now(),
			verloren = true,
			verloren_von_nutzer_id = $6,
			verloren_nutzer_name = $7
		WHERE medium_id = $8
	`

	name := nutzerModel.Vorname + " " + nutzerModel.Nachname
	_, err = d.DB.Exec(d.ctx, query, nil, nil, nil, MediumVerloren, nil, nutzerModel.NutzerId, name, payload.MediumId)
	if err != nil {
		log.Println("Error updating read-model:", err)
		return
	}

	if err := d.saveHistoryEvent(payload.MediumId, events.MediumVerlorenDurchBenutzerEventType, payload); err != nil {
		log.Println("Error saving history event:", err)
		return
	}
	d.notificationService.Notify(payload.MediumId)
}

func (d *DetailseiteProjection) saveHistoryEvent(mediumID string, eventType string, eventPayload any) error {
	payloadJSON, err := json.Marshal(eventPayload)
	if err != nil {
		return err
	}

	const query = `
		INSERT INTO medium_historie (medium_id, event_type, payload)
		VALUES ($1, $2, $3)
	`

	_, err = d.DB.Exec(d.ctx, query, mediumID, eventType, payloadJSON)
	return err
}
