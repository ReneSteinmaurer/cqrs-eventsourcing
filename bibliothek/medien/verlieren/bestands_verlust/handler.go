package bestands_verlust

import (
	"context"
	shared2 "cqrs-playground/bibliothek/medien/shared"
	"cqrs-playground/bibliothek/medien/verlieren/events"
	"cqrs-playground/shared"
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
	"time"
)

type MediumBestandsverlustHandler struct {
	eventStore   *shared.EventStore
	kafkaService *shared.KafkaService
	producer     sarama.SyncProducer
	ctx          context.Context
}

func NewMediumBestandsverlustHandler(
	ctx context.Context, eventStore *shared.EventStore, kafkaService *shared.KafkaService,
) *MediumBestandsverlustHandler {
	producer := kafkaService.NewSyncProducer()
	return &MediumBestandsverlustHandler{
		eventStore:   eventStore,
		kafkaService: kafkaService,
		producer:     producer,
		ctx:          ctx,
	}
}

func (v *MediumBestandsverlustHandler) Handle(cmd MediumBestandsverlustCommand) (string, error) {
	aggregateKey := cmd.MediumId
	aggregateType := shared2.MediumAggregateType

	if cmd.MediumId == "" {
		return "", errors.New("alle Felder muessen befuellt sein")
	}

	return aggregateKey, shared.RetryHandlerBasedOnVersionConflict(func() error {
		aggregateEvents, err := v.eventStore.GetEventsByAggregateId(v.ctx, aggregateKey, aggregateType)
		if err != nil {
			return err
		}

		aggregate := shared2.NewMediumAggregate(aggregateEvents)

		payload := events.NewMediumBestandsverlustEvent(cmd.MediumId, time.Now())
		err = aggregate.HandleMediumBestandsverlust(payload)
		if err != nil {
			return err
		}

		return v.SendEvent(payload, aggregateKey, aggregateType)
	})
}

func (v *MediumBestandsverlustHandler) SendEvent(payload events.MediumBestandsverlustEvent, aggregateKey, aggregateType string) error {
	eventType := events.MediumBestandsverlustEventType
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	version, err := v.eventStore.LoadCurrentVersion(v.ctx, aggregateKey, aggregateType)
	if err != nil {
		return err
	}

	event := shared.NewEvent(
		aggregateType,
		aggregateKey,
		eventType,
		version+1,
		payloadJSON)

	err = v.eventStore.Save(v.ctx, event)
	if err != nil {
		return err
	}
	err = v.kafkaService.SendEvent(v.producer, eventType, payloadJSON)
	if err != nil {
		panic(err)
	}
	return nil
}
