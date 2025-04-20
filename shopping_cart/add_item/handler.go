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
	EventStore *shared.EventStore
	ctx        context.Context
}

func NewAddItemHandler(ctx context.Context, eventStore *shared.EventStore) *AddItemHandler {
	return &AddItemHandler{
		EventStore: eventStore,
		ctx:        ctx,
	}
}

func (h *AddItemHandler) Handle(cmd AddItemToCartCommand) error {
	if cmd.Quantity <= 0 {
		return errors.New("quantity must be greater then zero")
	}

	eventPayload := ItemAddedToCartEvent{
		CartId:   cmd.CartId,
		Item:     cmd.Item,
		Quantity: cmd.Quantity,
	}

	payloadJSON, err := json.Marshal(eventPayload)
	if err != nil {
		return err
	}

	event := shared.Event{
		Id:        uuid.NewString(),
		Type:      "ItemAddedToCart",
		Timestamp: time.Now().UTC(),
		Payload:   payloadJSON,
	}

	return h.EventStore.Save(h.ctx, event)
}
