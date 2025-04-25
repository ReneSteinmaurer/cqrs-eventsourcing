package katalogisieren

import (
	"context"
	"cqrs-playground/bibliothek/medien/katalogisieren/events"
	shared2 "cqrs-playground/bibliothek/medien/shared"
	"cqrs-playground/shared"
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
)

type KatalogisiereMediumHandler struct {
	eventStore   *shared.EventStore
	kafkaService *shared.KafkaService
	producer     sarama.SyncProducer
	ctx          context.Context
}

func NewKatalogisiereMediumHandler(
	ctx context.Context, eventStore *shared.EventStore, kafkaService *shared.KafkaService,
) *KatalogisiereMediumHandler {
	producer := kafkaService.NewSyncProducer()
	return &KatalogisiereMediumHandler{
		eventStore:   eventStore,
		kafkaService: kafkaService,
		producer:     producer,
		ctx:          ctx,
	}
}

func (k *KatalogisiereMediumHandler) Handle(cmd KatalogisiereMediumCommand) (string, error) {
	aggregateKey := cmd.MediumId
	aggregateType := shared2.MediumAggregateType

	if cmd.MediumId == "" || cmd.Signature == "" {
		return "", errors.New("alle Felder müssen befüllt sein")
	}

	return aggregateKey, shared.RetryHandlerBasedOnVersionConflict(func() error {
		aggregateEvents, err := k.eventStore.GetEventsByAggregateId(k.ctx, aggregateKey, aggregateType)
		if err != nil {
			return err
		}

		aggregate := shared2.NewMediumAggregate(aggregateEvents)

		payload := events.NewMediumKatalogisiertEvent(cmd.MediumId, cmd.Signature, cmd.Standort, cmd.ExemplarCode)
		err = aggregate.HandleMediumKatalogisieren(payload)
		if err != nil {
			return err
		}

		return k.SendEvent(payload, aggregateKey, aggregateType)
	})
}

func (k *KatalogisiereMediumHandler) SendEvent(payload events.MediumKatalogisiertEvent, aggregateKey, aggregateType string) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	version, err := k.eventStore.LoadCurrentVersion(k.ctx, aggregateKey, aggregateType)
	if err != nil {
		return err
	}

	event := shared.NewEvent(
		aggregateType,
		aggregateKey,
		events.MediumKatalogisiertEventType,
		version+1,
		payloadJSON)

	err = k.eventStore.Save(k.ctx, event)
	if err != nil {
		return err
	}
	err = k.kafkaService.SendEvent(k.producer, events.MediumKatalogisiertEventType, payloadJSON)
	if err != nil {
		panic(err)
	}
	return nil
}
