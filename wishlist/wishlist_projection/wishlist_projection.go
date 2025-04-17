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
	projectionUpdater *shared.ProjectionStateUpdater
	ctx               context.Context
	cancel            context.CancelFunc
}

func NewWishlistProjection(ctx context.Context, eventStore *shared.EventStore, db *pgxpool.Pool, projectionStateUpdater *shared.ProjectionStateUpdater) *WishlistProjection {
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
	}
}

func (cp *WishlistProjection) Start(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cp.updateProjection()
		case <-cp.ctx.Done():
			log.Println("WishlistProjection beendet.")
			return
		}
	}
}

func (cp *WishlistProjection) updateProjection() {
	events, err := cp.EventStore.GetEventsSince(cp.ctx, cp.LastUpdate.UTC())
	if err != nil {
		log.Println("Error fetching events:", err)
		return
	}

	for _, event := range events {
		if event.Type == add_item.ItemAddedToWishlistEventTypeV1 {
			cp.applyItemAddedV1(event)
		} else if event.Type == add_item.ItemAddedToWishlistEventTypeV2 {
			cp.applyItemAddedV2(event)
		}

		now := time.Now().UTC()
		cp.LastUpdate = now
		cp.projectionUpdater.UpdateProjectionState(projectionName, event.Id, now)
	}
}

func (cp *WishlistProjection) applyItemAddedV1(event shared.Event) {
	log.Println("Item added to wishlist")
	var payload add_item.ItemAddedToWishlistEventV1
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
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

func (cp *WishlistProjection) applyItemAddedV2(event shared.Event) {
	log.Println("Item added to wishlist")
	var payload add_item.ItemAddedToWishlistEventV2
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
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
