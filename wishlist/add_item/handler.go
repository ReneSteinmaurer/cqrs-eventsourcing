package add_item

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

type AddItemHandler struct {
	eventStore   *shared.EventStore
	kafkaService *shared.KafkaService
	producer     sarama.SyncProducer
	ctx          context.Context
}

func NewAddItemHandler(
	ctx context.Context, eventStore *shared.EventStore, kafkaService *shared.KafkaService,
) *AddItemHandler {
	producer := kafkaService.NewSyncProducer()
	return &AddItemHandler{
		eventStore:   eventStore,
		kafkaService: kafkaService,
		producer:     producer,
		ctx:          ctx,
	}
}

func (a *AddItemHandler) HandleV1(cmd AddItemToWishlistCommandV1) error {
	if cmd.WishlistId < 0 {
		return errors.New("wishlist id cannot be negative")
	}
	if cmd.Item == "" {
		return errors.New("item cannot be empty string")
	}

	payload := events.ItemAddedToWishlistEventV1{
		WishlistId: cmd.WishlistId,
		Item:       cmd.Item,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	version, err := a.eventStore.LoadCurrentVersion(a.ctx, strconv.Itoa(cmd.WishlistId), wishlist.AggregateType)
	if err != nil {
		panic(err)
	}

	event := shared.NewEvent(
		wishlist.AggregateType,
		strconv.Itoa(cmd.WishlistId),
		events.ItemAddedToWishlistEventTypeV1,
		version+1, payloadJSON)

	err = a.eventStore.Save(a.ctx, event)
	if err != nil {
		return err
	}

	err = a.kafkaService.SendEvent(a.producer, events.ItemAddedToWishlistEventTypeV1, payloadJSON)
	if err != nil {
		panic(err)
	}
	return nil
}

func (a *AddItemHandler) HandleV2(cmd AddItemToWishlistCommandV2) error {
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

		payload := events.ItemAddedToWishlistEventV2{
			WishlistId: cmd.WishlistId,
			Item:       cmd.Item,
			UserId:     cmd.UserId,
		}

		err = wishlistAggregate.HandleAddItem(payload)
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
			events.ItemAddedToWishlistEventTypeV2,
			version+1,
			payloadJSON)

		err = a.eventStore.Save(a.ctx, event)
		if err != nil {
			return err
		}
		err = a.kafkaService.SendEvent(a.producer, events.ItemAddedToWishlistEventTypeV2, payloadJSON)
		if err != nil {
			panic(err)
		}
		return nil
	})
}
