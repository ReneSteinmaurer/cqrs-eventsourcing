package verloren_duch_benutzer

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

type MediumVerlorenDurchBenutzerHandler struct {
	eventStore   *shared.EventStore
	kafkaService *shared.KafkaService
	producer     sarama.SyncProducer
	ctx          context.Context
}

func NewMediumVerlorenDurchBenutzerHandler(
	ctx context.Context, eventStore *shared.EventStore, kafkaService *shared.KafkaService,
) *MediumVerlorenDurchBenutzerHandler {
	producer := kafkaService.NewSyncProducer()
	return &MediumVerlorenDurchBenutzerHandler{
		eventStore:   eventStore,
		kafkaService: kafkaService,
		producer:     producer,
		ctx:          ctx,
	}
}

func (v *MediumVerlorenDurchBenutzerHandler) Handle(cmd MediumVerlorenDurchBenutzerCommand) (string, error) {
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

		payload := events.NewMediumVerlorenEvent(cmd.MediumId, cmd.NutzerId, time.Now())
		err = aggregate.HandleMediumVerlorenDurchBenutzer(payload)
		if err != nil {
			return err
		}

		return v.SendEvent(payload, aggregateKey, aggregateType)
	})
}

func (v *MediumVerlorenDurchBenutzerHandler) SendEvent(payload events.MediumVerlorenDurchBenutzerEvent, aggregateKey, aggregateType string) error {
	eventType := events.MediumVerlorenDurchBenutzerEventType
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
