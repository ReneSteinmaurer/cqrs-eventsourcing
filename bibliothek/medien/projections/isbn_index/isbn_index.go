package isbn_index

import (
	"context"
	shared2 "cqrs-playground/bibliothek/medien/erwerben/events"
	"cqrs-playground/shared"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

const projectionName = "IsbnIndexProjection"

type ISBNIndexProjection struct {
	EventStore        *shared.EventStore
	DB                *pgxpool.Pool
	KafkaService      *shared.KafkaService
	projectionUpdater *shared.ProjectionStateUpdater
	ctx               context.Context
	cancel            context.CancelFunc
}

func NewISBNIndexProjection(
	ctx context.Context, eventStore *shared.EventStore, db *pgxpool.Pool, projectionStateUpdater *shared.ProjectionStateUpdater,
	kafkaService *shared.KafkaService,
) *ISBNIndexProjection {
	ctx, cancel := context.WithCancel(ctx)
	return &ISBNIndexProjection{
		ctx:               ctx,
		cancel:            cancel,
		projectionUpdater: projectionStateUpdater,
		EventStore:        eventStore,
		DB:                db,
		KafkaService:      kafkaService,
	}
}

func (cp *ISBNIndexProjection) listenToEvent(eventType string, applyFunc func([]byte)) {
	consumer := cp.KafkaService.NewConsumerOffsetNewest(eventType)
	defer func() {
		log.Printf("Closing consumer for %s...\n", eventType)
		_ = consumer.Close()
	}()

	msgs := consumer.Messages()

	for {
		select {
		case <-cp.ctx.Done():
			log.Printf("%s listener stopped\n", eventType)
			return
		case msg := <-msgs:
			applyFunc(msg.Value)
		}
	}
}

func (cp *ISBNIndexProjection) Start() {
	go cp.listenToEvent(shared2.MediumErworbenEventType, cp.applyMediumErworben)
}

func (cp *ISBNIndexProjection) applyMediumErworben(payloadJSON []byte) {
	var payload shared2.MediumErworbenEvent
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		log.Println("Error unmarshalling event:", err)
		return
	}

	const query = `
		insert into isbn_index (isbn, medium_id)
        values ($1, $2)
		on conflict do nothing 
	`

	_, err := cp.DB.Exec(cp.ctx, query, payload.ISBN, payload.MediumId)
	if err != nil {
		log.Println("Error updating read-model:", err)
	}
}
