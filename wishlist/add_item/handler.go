package add_item

import (
	"context"
	"cqrs-playground/shared"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"time"
)

type AddItemHandler struct {
	eventStore *shared.EventStore
	ctx        context.Context
}

func NewAddItemHandler(ctx context.Context, eventStore *shared.EventStore) *AddItemHandler {
	return &AddItemHandler{
		eventStore: eventStore,
		ctx:        ctx,
	}
}

func (a *AddItemHandler) HandleV1(cmd AddItemToWishlistCommandV1) error {
	if cmd.WishlistId < 0 {
		return errors.New("wishlist id cannot be negative")
	}
	if cmd.Item == "" {
		return errors.New("item cannot be empty string")
	}

	payload := ItemAddedToWishlistEventV1{
		WishlistId: cmd.WishlistId,
		Item:       cmd.Item,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	event := shared.Event{
		Id:        uuid.NewString(),
		Type:      ItemAddedToWishlistEventTypeV1,
		Timestamp: time.Now().UTC(),
		Payload:   payloadJSON,
	}

	return a.eventStore.Save(a.ctx, event)
}

func (a *AddItemHandler) HandleV2(cmd AddItemToWishlistCommandV2) error {
	if cmd.WishlistId < 0 {
		return errors.New("wishlist id cannot be negative")
	}
	if cmd.Item == "" {
		return errors.New("item cannot be empty string")
	}
	if cmd.UserId == "" {
		return errors.New("userId cannot be empty string")
	}

	payload := ItemAddedToWishlistEventV2{
		WishlistId: cmd.WishlistId,
		Item:       cmd.Item,
		UserId:     cmd.UserId,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	event := shared.Event{
		Id:        uuid.NewString(),
		Type:      ItemAddedToWishlistEventTypeV2,
		Timestamp: time.Now().UTC(),
		Payload:   payloadJSON,
	}

	return a.eventStore.Save(a.ctx, event)
}
