package remove_item

import (
	"context"
	"cqrs-playground/shared"
	"cqrs-playground/wishlist"
	"cqrs-playground/wishlist/events"
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
	"strconv"
)

type RemoveItemHandler struct {
	eventStore   *shared.EventStore
	kafkaService *shared.KafkaService
	producer     sarama.SyncProducer
	ctx          context.Context
}

func NewRemoveItemHandler(
	ctx context.Context, eventStore *shared.EventStore, kafkaService *shared.KafkaService,
) *RemoveItemHandler {
	producer := kafkaService.NewSyncProducer()
	return &RemoveItemHandler{
		eventStore:   eventStore,
		kafkaService: kafkaService,
		producer:     producer,
		ctx:          ctx,
	}
}

func (a *RemoveItemHandler) HandleV1(cmd RemoveItemFromWishlistCommandV1) error {
	var aggregateKey = strconv.Itoa(cmd.WishlistId)
	var aggregateType = wishlist.AggregateType
	if cmd.WishlistId < 0 {
		return errors.New("wishlist id cannot be negative")
	}
	if cmd.Item == "" {
		return errors.New("item cannot be empty string")
	}
	if cmd.UserId == "" {
		return errors.New("userId cannot be empty string")
	}

	return shared.RetryHandlerLogic(func() error {
		aggregateEvents, err := a.eventStore.GetEventsByAggregateId(a.ctx, aggregateKey, aggregateType)
		if err != nil {
			return err
		}

		wishlistAggregate := wishlist.NewWishlistAggregateFrom(aggregateEvents)

		version, err := a.eventStore.LoadCurrentVersion(a.ctx, aggregateKey, aggregateType)
		if err != nil {
			return err
		}

		payload := events.ItemRemovedFromWishlistV1{
			WishlistId: cmd.WishlistId,
			Item:       cmd.Item,
			UserId:     cmd.UserId,
		}

		err = wishlistAggregate.HandleRemoveItem(payload)
		if err != nil {
			return err
		}

		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		event := shared.NewEvent(
			aggregateType,
			aggregateKey,
			events.ItemRemovedFromWishlistTypeV1,
			version+1,
			payloadJSON)

		err = a.eventStore.Save(a.ctx, event)
		if err != nil {
			return err
		}
		err = a.kafkaService.SendEvent(a.producer, events.ItemRemovedFromWishlistTypeV1, payloadJSON)
		if err != nil {
			panic(err)
		}
		return nil
	})
}
