package bestands_verlust

import (
	"context"
	shared2 "cqrs-playground/bibliothek/medien/shared"
	"cqrs-playground/bibliothek/medien/wiederaufgefunden/events"
	"cqrs-playground/shared"
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
	"time"
)

type MediumWiederaufgefundenDurchNutzerHandler struct {
	eventStore   *shared.EventStore
	kafkaService *shared.KafkaService
	producer     sarama.SyncProducer
	ctx          context.Context
}

func NewMediumWiederaufgefundenDurchNutzerHandler(
	ctx context.Context, eventStore *shared.EventStore, kafkaService *shared.KafkaService,
) *MediumWiederaufgefundenDurchNutzerHandler {
	producer := kafkaService.NewSyncProducer()
	return &MediumWiederaufgefundenDurchNutzerHandler{
		eventStore:   eventStore,
		kafkaService: kafkaService,
		producer:     producer,
		ctx:          ctx,
	}
}

func (v *MediumWiederaufgefundenDurchNutzerHandler) Handle(cmd MediumWiederaufgefundenDurchNutzerCommand) (string, error) {
	aggregateKey := cmd.MediumId
	aggregateType := shared2.MediumAggregateType

	if cmd.MediumId == "" || cmd.NutzerId == "" {
		return "", errors.New("alle Felder muessen befuellt sein")
	}

	return aggregateKey, shared.RetryHandlerBasedOnVersionConflict(func() error {
		aggregateEvents, err := v.eventStore.GetEventsByAggregateId(v.ctx, aggregateKey, aggregateType)
		if err != nil {
			return err
		}

		aggregate := shared2.NewMediumAggregate(aggregateEvents)

		payload := events.NewMediumWiederaufgefundenDurchNutzerEvent(cmd.MediumId, cmd.NutzerId, time.Now())
		err = aggregate.HandleMediumWiederaufgefundenDurchNutzer(payload)
		if err != nil {
			return err
		}

		return v.SendEvent(payload, aggregateKey, aggregateType)
	})
}

func (v *MediumWiederaufgefundenDurchNutzerHandler) SendEvent(payload events.MediumWiederaufgefundenDurchNutzerEvent, aggregateKey, aggregateType string) error {
	eventType := events.MediumWiederaufgefundenDurchNutzerEventType
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
