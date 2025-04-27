package nutzer

import (
	"context"
	"cqrs-playground/bibliothek/nutzer/registrierung/events"
	"cqrs-playground/shared"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type NutzerProjection struct {
	EventStore   *shared.EventStore
	DB           *pgxpool.Pool
	KafkaService *shared.KafkaService
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewNutzerProjection(
	ctx context.Context, eventStore *shared.EventStore, db *pgxpool.Pool,
	kafkaService *shared.KafkaService,
) *NutzerProjection {
	ctx, cancel := context.WithCancel(ctx)
	return &NutzerProjection{
		ctx:          ctx,
		cancel:       cancel,
		EventStore:   eventStore,
		DB:           db,
		KafkaService: kafkaService,
	}
}

func (n *NutzerProjection) Start() {
	go shared.ListenToEvent(n.ctx, n.KafkaService, events.NutzerRegistriertEventType, n.applyNutzerRegistriert)
}

func (n *NutzerProjection) applyNutzerRegistriert(payloadJSON []byte) {
	var payload events.NutzerRegistriertEvent
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		log.Println("Error unmarshalling event:", err)
		return
	}

	const query = `
		insert into nutzer (nutzer_id, email, vorname, nachname)
        values ($1, $2, $3, $4)
		on conflict do nothing 
	`

	_, err := n.DB.Exec(n.ctx, query, payload.NutzerId, payload.Email, payload.Vorname, payload.Nachname)
	if err != nil {
		log.Println("Error updating read-model:", err)
	}
}
