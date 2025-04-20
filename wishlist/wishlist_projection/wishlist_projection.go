package wishlist_projection

import (
	"context"
	"cqrs-playground/shared"
	"cqrs-playground/wishlist/events"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

const projectionName = "WishlistProjection"

type WishlistProjection struct {
	EventStore        *shared.EventStore
	DB                *pgxpool.Pool
	KafkaService      *shared.KafkaService
	projectionUpdater *shared.ProjectionStateUpdater
	ctx               context.Context
	cancel            context.CancelFunc
}

func NewWishlistProjection(
	ctx context.Context, eventStore *shared.EventStore, db *pgxpool.Pool, projectionStateUpdater *shared.ProjectionStateUpdater,
	kafkaService *shared.KafkaService,
) *WishlistProjection {
	ctx, cancel := context.WithCancel(ctx)
	return &WishlistProjection{
		ctx:               ctx,
		cancel:            cancel,
		projectionUpdater: projectionStateUpdater,
		EventStore:        eventStore,
		DB:                db,
		KafkaService:      kafkaService,
	}
}

func (cp *WishlistProjection) Start() {
	go cp.listenToEvent(events.ItemAddedToWishlistEventTypeV1, cp.applyItemAddedV1)
	go cp.listenToEvent(events.ItemAddedToWishlistEventTypeV2, cp.applyItemAddedV2)
	go cp.listenToEvent(events.ItemRemovedFromWishlistTypeV1, cp.applyItemRemovedV1)
}

func (cp *WishlistProjection) listenToEvent(eventType string, applyFunc func([]byte)) {
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

func (cp *WishlistProjection) applyItemAddedV1(payloadJSON []byte) {
	log.Println("Item added to wishlist v1")
	var payload events.ItemAddedToWishlistEventV1
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		log.Println("Error unmarshalling event:", err)
		return
	}

	const query = `
		insert into wishlist_items (wishlist_id, item)
        values ($1, $2)
		on conflict do nothing 
	`

	_, err := cp.DB.Exec(cp.ctx, query, payload.WishlistId, payload.Item)
	if err != nil {
		log.Println("Error updating read-model:", err)
	}
}

func (cp *WishlistProjection) applyItemAddedV2(payloadJSON []byte) {
	log.Println("Item added to wishlist v2")
	var payload events.ItemAddedToWishlistEventV2
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		log.Println("Error unmarshalling event:", err)
		return
	}

	const query = `
		insert into wishlist_items (wishlist_id, item, user_id)
        values ($1, $2, $3)
		on conflict do nothing 
	`

	_, err := cp.DB.Exec(cp.ctx, query, payload.WishlistId, payload.Item, payload.UserId)
	if err != nil {
		log.Println("Error updating read-model:", err)
	}
}

func (cp *WishlistProjection) applyItemRemovedV1(payloadJSON []byte) {
	log.Println("Item removed from wishlist v1")
	var payload events.ItemRemovedFromWishlistV1
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		log.Println("Error unmarshalling event:", err)
		return
	}

	const query = `
		delete from wishlist_items
		where wishlist_id = $1 and item = $2 and user_id = $3
	`

	_, err := cp.DB.Exec(cp.ctx, query, payload.WishlistId, payload.Item, payload.UserId)
	if err != nil {
		log.Println("Error updating read-model:", err)
	}
}
