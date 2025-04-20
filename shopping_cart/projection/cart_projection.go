package projection

import (
	"context"
	"cqrs-playground/shared"
	"cqrs-playground/shopping_cart/add_item"
	"cqrs-playground/shopping_cart/remove_item"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

const projectionName = "CartProjection"

type CartProjection struct {
	EventStore        *shared.EventStore
	DB                *pgxpool.Pool
	KafkaService      *shared.KafkaService
	LastUpdate        time.Time
	projectionUpdater *shared.ProjectionStateUpdater
	ctx               context.Context
	cancel            context.CancelFunc
}

func NewCartProjection(
	ctx context.Context, eventStore *shared.EventStore, db *pgxpool.Pool, projectionStateUpdater *shared.ProjectionStateUpdater,
	kafkaService *shared.KafkaService,
) *CartProjection {
	ctx, cancel := context.WithCancel(ctx)
	projectionStatus, err := eventStore.GetLastUpdateFromProjection(ctx, projectionName)
	if err != nil {
		panic(err)
	}
	return &CartProjection{
		ctx:               ctx,
		cancel:            cancel,
		projectionUpdater: projectionStateUpdater,
		EventStore:        eventStore,
		DB:                db,
		KafkaService:      kafkaService,
		LastUpdate:        projectionStatus.LastProcessedTimestamp,
	}
}

func (cp *CartProjection) Start(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cp.updateProjection()
		case <-cp.ctx.Done():
			log.Println("CartProjection beendet.")
			return
		}
	}
}

func (cp *CartProjection) updateProjection() {
	events, err := cp.EventStore.GetEventsSince(cp.ctx, cp.LastUpdate.UTC())
	if err != nil {
		log.Println("Error fetching events:", err)
		return
	}

	for _, event := range events {
		if event.Type == "ItemAddedToCart" {
			cp.applyItemAdded(event)
		} else if event.Type == "ItemRemovedFromCart" {
			cp.applyItemRemoved(event)
		}

		now := time.Now().UTC()
		cp.LastUpdate = now
		cp.projectionUpdater.UpdateProjectionState(projectionName, event.Id, now)
	}
}

func (cp *CartProjection) applyItemAdded(event shared.Event) {
	log.Println("Item added to cart")
	var payload add_item.ItemAddedToCartEvent
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		log.Println("Error unmarshalling event:", err)
		return
	}

	const query = `
		insert into cart_items (cart_id, item, quantity)
        values ($1, $2, $3)
        on conflict (cart_id, item)
        do update set quantity = cart_items.quantity + EXCLUDED.quantity
	`

	_, err := cp.DB.Exec(cp.ctx, query, payload.CartId, payload.Item, payload.Quantity)
	if err != nil {
		log.Println("Error updating read-model:", err)
	}
}

func (cp *CartProjection) applyItemRemoved(event shared.Event) {
	log.Println("Item removed from cart")
	var payload remove_item.ItemRemovedFromCartEvent
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		log.Println("Error unmarshalling event:", err)
		return
	}

	const query = `
		update cart_items 
		set quantity = quantity-1 
		where item = $1 and cart_id = $2
		and quantity > 0
	`

	_, err := cp.DB.Exec(cp.ctx, query, payload.Item, payload.CartId)
	if err != nil {
		log.Println("Error updating read-model:", err)
	}
}

func (cp *CartProjection) Stop() {
	cp.cancel()
}
