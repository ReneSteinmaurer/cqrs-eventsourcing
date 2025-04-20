package remove_item

import (
	"context"
	"cqrs-playground/shared"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"time"
)

type RemoveItemHandler struct {
	EventStore *shared.EventStore
	ctx        context.Context
}

func NewRemoveItemHandler(ctx context.Context, eventStore *shared.EventStore) *RemoveItemHandler {
	return &RemoveItemHandler{
		EventStore: eventStore,
		ctx:        ctx,
	}
}

func (h *RemoveItemHandler) Handle(cmd RemoveItemFromCartCommand) error {
	if cmd.CartId < 0 {
		return errors.New("cartId cannot be negative")
	}

	payload := ItemRemovedFromCartEvent{
		CartId: cmd.CartId,
		Item:   cmd.Item,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	event := shared.Event{
		Id:        uuid.NewString(),
		Type:      "ItemRemovedFromCart",
		Timestamp: time.Now().UTC(),
		Payload:   payloadJSON,
	}

	return h.EventStore.Save(h.ctx, event)
}
