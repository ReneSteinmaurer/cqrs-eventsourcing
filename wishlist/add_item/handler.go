package add_item

import (
	"context"
	"cqrs-playground/shared"
	"cqrs-playground/wishlist"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
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

	payload := ItemAddedToWishlistEventV1{
		WishlistId: cmd.WishlistId,
		Item:       cmd.Item,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	version, err := a.eventStore.LoadCurrentVersion(a.ctx, wishlist.AggregateId)
	if err != nil {
		panic(err)
	}

	event := shared.NewEvent(wishlist.AggregateId, ItemAddedToWishlistEventTypeV1, version+1, payloadJSON)
	err = a.eventStore.Save(a.ctx, event)
	if err != nil {
		return err
	}

	err = a.kafkaService.SendEvent(a.producer, ItemAddedToWishlistEventTypeV1, payloadJSON)
	if err != nil {
		panic(err)
	}
	return nil
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

	return shared.RetryHandlerLogic(func() error {
		aggregateId := aggregateIdForWishlist(cmd.WishlistId)
		version, err := a.eventStore.LoadCurrentVersion(a.ctx, aggregateId)
		if err != nil {
			return err
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

		event := shared.NewEvent(aggregateId, ItemAddedToWishlistEventTypeV2, version+1, payloadJSON)
		err = a.eventStore.Save(a.ctx, event)
		if err != nil {
			return err
		}
		err = a.kafkaService.SendEvent(a.producer, ItemAddedToWishlistEventTypeV2, payloadJSON)
		if err != nil {
			panic(err)
		}
		return nil
	})
}

func aggregateIdForWishlist(id int) string {
	return fmt.Sprintf("%s-%d", wishlist.AggregateId, id)
}
