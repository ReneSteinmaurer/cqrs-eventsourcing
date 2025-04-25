package verliehen

import (
	"context"
	"cqrs-playground/bibliothek/medien/ausleihen/events"
	shared2 "cqrs-playground/bibliothek/medien/rueckgeben/events"
	"cqrs-playground/shared"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type MediumVerliehenProjection struct {
	EventStore   *shared.EventStore
	DB           *pgxpool.Pool
	KafkaService *shared.KafkaService
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewMediumVerliehenProjection(
	ctx context.Context, eventStore *shared.EventStore, db *pgxpool.Pool,
	kafkaService *shared.KafkaService,
) *MediumVerliehenProjection {
	ctx, cancel := context.WithCancel(ctx)
	return &MediumVerliehenProjection{
		ctx:          ctx,
		cancel:       cancel,
		EventStore:   eventStore,
		DB:           db,
		KafkaService: kafkaService,
	}
}

func (mv *MediumVerliehenProjection) listenToEvent(eventType string, applyFunc func([]byte)) {
	consumer := mv.KafkaService.NewConsumerOffsetNewest(eventType)
	defer func() {
		log.Printf("Closing consumer for %s...\n", eventType)
		_ = consumer.Close()
	}()

	msgs := consumer.Messages()

	for {
		select {
		case <-mv.ctx.Done():
			log.Printf("%s listener stopped\n", eventType)
			return
		case msg := <-msgs:
			applyFunc(msg.Value)
		}
	}
}

func (mv *MediumVerliehenProjection) Start() {
	go mv.listenToEvent(events.MediumVerliehenEventType, mv.applyMediumVerliehen)
	go mv.listenToEvent(shared2.MediumZurueckgegebenEventType, mv.applyMediumZurueckgegeben)
}

func (mv *MediumVerliehenProjection) applyMediumVerliehen(payloadJSON []byte) {
	var payload events.MediumVerliehenEvent
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		log.Println("Error unmarshalling event:", err)
		return
	}

	const query = `
		insert into medium_verliehen (medium_id, verliehen_von, verliehen_bis, nutzer_id)
        values ($1, $2, $3, $4)
		on conflict do nothing 
	`

	_, err := mv.DB.Exec(mv.ctx, query, payload.MediumId, payload.Von, payload.Bis, payload.NutzerId)
	if err != nil {
		log.Println("Error updating read-model:", err)
	}
}

func (mv *MediumVerliehenProjection) applyMediumZurueckgegeben(payloadJSON []byte) {
	var payload shared2.MediumZurueckgegebenEvent
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		log.Println("Error unmarshalling event:", err)
		return
	}

	const query = `
		delete from medium_verliehen where medium_id = $1 and nutzer_id = $2;
	`

	_, err := mv.DB.Exec(mv.ctx, query, payload.MediumId, payload.NutzerId)
	if err != nil {
		log.Println("Error updating read-model:", err)
	}
}
