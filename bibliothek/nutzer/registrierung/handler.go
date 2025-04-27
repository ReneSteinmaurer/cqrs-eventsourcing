package registrierung

import (
	"context"
	"cqrs-playground/bibliothek/nutzer"
	"cqrs-playground/bibliothek/nutzer/registrierung/events"
	"cqrs-playground/shared"
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
)

type NutzerRegistrierungHandler struct {
	eventStore   *shared.EventStore
	kafkaService *shared.KafkaService
	producer     sarama.SyncProducer
	ctx          context.Context
}

func NewNutzerRegistrierungHandler(
	ctx context.Context, eventStore *shared.EventStore, kafkaService *shared.KafkaService,
) *NutzerRegistrierungHandler {
	producer := kafkaService.NewSyncProducer()
	return &NutzerRegistrierungHandler{
		eventStore:   eventStore,
		kafkaService: kafkaService,
		producer:     producer,
		ctx:          ctx,
	}
}

func (n *NutzerRegistrierungHandler) Handle(cmd NutzerRegistrierungCommand) error {
	aggregateKey := cmd.Email
	aggregateType := nutzer.AggregateTypeRegistrierung
	if cmd.Email == "" {
		return errors.New("die Email wird benötigt")
	}

	return shared.RetryHandlerBasedOnVersionConflict(func() error {
		aggregateEvents, err := n.eventStore.GetEventsByAggregateId(n.ctx, aggregateKey, aggregateType)
		if err != nil {
			return err
		}

		aggregate := nutzer.NewNutzerAggregate(aggregateEvents)

		version, err := n.eventStore.LoadCurrentVersion(n.ctx, aggregateKey, aggregateType)
		if err != nil {
			return err
		}

		payload := events.NewNutzerRegistrierungEvent(cmd.Email, cmd.Vorname, cmd.Nachname, aggregateKey)

		err = aggregate.HandleRegistriereNutzer(payload)
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
			events.NutzerRegistriertEventType,
			version+1,
			payloadJSON)

		err = n.eventStore.Save(n.ctx, event)
		if err != nil {
			return err
		}
		err = n.kafkaService.SendEvent(n.producer, events.NutzerRegistriertEventType, payloadJSON)
		if err != nil {
			panic(err)
		}
		return nil
	})
}
