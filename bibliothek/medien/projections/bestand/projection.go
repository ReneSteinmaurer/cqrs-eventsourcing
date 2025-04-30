package bestand

import (
	"context"
	"cqrs-playground/bibliothek/medien/erwerben/events"
	shared2 "cqrs-playground/bibliothek/medien/katalogisieren/events"
	"cqrs-playground/shared"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type MedienProjection struct {
	EventStore          *shared.EventStore
	DB                  *pgxpool.Pool
	KafkaService        *shared.KafkaService
	NotificationService *shared.NotificationService
	ctx                 context.Context
	cancel              context.CancelFunc
}

func NewMediumBestandProjection(
	ctx context.Context, eventStore *shared.EventStore, db *pgxpool.Pool,
	kafkaService *shared.KafkaService, notificationService *shared.NotificationService,
) *MedienProjection {
	ctx, cancel := context.WithCancel(ctx)
	return &MedienProjection{
		ctx:                 ctx,
		cancel:              cancel,
		EventStore:          eventStore,
		DB:                  db,
		KafkaService:        kafkaService,
		NotificationService: notificationService,
	}
}

func (cp *MedienProjection) Start() {
	go shared.ListenToEvent(cp.ctx, cp.KafkaService, events.MediumErworbenEventType, cp.applyMediumErworben)
	go shared.ListenToEvent(cp.ctx, cp.KafkaService, shared2.MediumKatalogisiertEventType, cp.applyMediumKatalogisiert)
}

func (cp *MedienProjection) applyMediumErworben(payloadJSON []byte) {
	var payload events.MediumErworbenEvent
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		log.Println("Error unmarshalling event:", err)
		return
	}

	const query = `
		insert into medium_bestand (medium_id, isbn, medium_type, name, genre)
        values ($1, $2, $3, $4, $5)
		on conflict do nothing 
	`

	_, err := cp.DB.Exec(cp.ctx, query, payload.MediumId, payload.ISBN, payload.MediumType, payload.Name, payload.Genre)
	if err != nil {
		log.Println("Error updating read-model:", err)
	}
	cp.NotificationService.NotifyProjectionUpdated(payload.MediumId)
}

func (cp *MedienProjection) applyMediumKatalogisiert(payloadJSON []byte) {
	var payload shared2.MediumKatalogisiertEvent
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		log.Println("Error unmarshalling event:", err)
		return
	}

	const query = `
		update medium_bestand
		set
			signature = $2,
			standort = $3,
			exemplar_code = $4,
		    katalogisiert = true
		where medium_id = $1
	`

	_, err := cp.DB.Exec(cp.ctx, query,
		payload.MediumId,
		payload.Signature,
		payload.Standort,
		payload.ExemplarCode,
	)
	if err != nil {
		log.Println("Error updating read-model:", err)
		return
	}
	cp.NotificationService.NotifyProjectionUpdated(payload.MediumId)
}
