package wishlist_projection

import (
	"context"
	"cqrs-playground/shared"
	"cqrs-playground/wishlist/add_item"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

const projectionName = "WishlistProjection"

type WishlistProjection struct {
	EventStore        *shared.EventStore
	DB                *pgxpool.Pool
	LastUpdate        time.Time
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
	projectionStatus, err := eventStore.GetLastUpdateFromProjection(ctx, projectionName)
	if err != nil {
		panic(err)
	}

	return &WishlistProjection{
		ctx:               ctx,
		cancel:            cancel,
		projectionUpdater: projectionStateUpdater,
		EventStore:        eventStore,
		DB:                db,
		LastUpdate:        projectionStatus.LastProcessedTimestamp,
		KafkaService:      kafkaService,
	}
}

func (cp *WishlistProjection) Start() {
	go cp.itemAddedToWishlistEventType1Listener()
	go cp.itemAddedToWishlistEventType2Listener()
}

func (cp *WishlistProjection) itemAddedToWishlistEventType1Listener() {
	consumer := cp.KafkaService.NewConsumerOffsetNewest(add_item.ItemAddedToWishlistEventTypeV1)
	defer func() {
		log.Println("Closing V1 consumer...")
		_ = consumer.Close()
	}()

	msgs := consumer.Messages()

	for {
		select {
		case <-cp.ctx.Done():
			log.Println("itemAddedToWishlistEventType1Listener stopped")
			return
		case msg := <-msgs:
			cp.applyItemAddedV1(msg.Value)
		}
	}
}

func (cp *WishlistProjection) itemAddedToWishlistEventType2Listener() {
	consumer := cp.KafkaService.NewConsumerOffsetNewest(add_item.ItemAddedToWishlistEventTypeV2)
	defer func() {
		log.Println("Closing V2 consumer...")
		_ = consumer.Close()
	}()

	msgs := consumer.Messages()

	for {
		select {
		case <-cp.ctx.Done():
			log.Println("itemAddedToWishlistEventType2Listener stopped")
			return
		case msg := <-msgs:
			cp.applyItemAddedV2(msg.Value)
		}
	}
}

func (cp *WishlistProjection) applyItemAddedV1(payloadJSON []byte) {
	log.Println("Item added to wishlist v1")
	var payload add_item.ItemAddedToWishlistEventV1
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
	var payload add_item.ItemAddedToWishlistEventV2
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
