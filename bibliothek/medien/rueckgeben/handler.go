package rueckgeben

import (
	"context"
	shared2 "cqrs-playground/bibliothek/medien/shared"
	"cqrs-playground/bibliothek/medien/verliehen_projection"
	"cqrs-playground/shared"
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type MediumRueckgabeHandler struct {
	eventStore      *shared.EventStore
	kafkaService    *shared.KafkaService
	producer        sarama.SyncProducer
	verliehenReader *verliehen_projection.MediumVerliehenReader
	ctx             context.Context
}

func NewMediumRueckgabeHandler(
	ctx context.Context, eventStore *shared.EventStore, kafkaService *shared.KafkaService, db *pgxpool.Pool,
) *MediumRueckgabeHandler {
	producer := kafkaService.NewSyncProducer()
	vr := verliehen_projection.NewMediumVerliehenReader(db)
	return &MediumRueckgabeHandler{
		eventStore:      eventStore,
		kafkaService:    kafkaService,
		producer:        producer,
		verliehenReader: vr,
		ctx:             ctx,
	}
}

func (v *MediumRueckgabeHandler) Handle(cmd MediumRueckgebenCommand) error {
	aggregateKey := cmd.MediumId
	aggregateType := shared2.MediumAggregateType

	if cmd.MediumId == "" || cmd.NutzerId == "" {
		return errors.New("alle Felder muessen befuellt sein")
	}

	return shared.RetryHandlerLogic(func() error {
		aggregateEvents, err := v.eventStore.GetEventsByAggregateId(v.ctx, aggregateKey, aggregateType)
		if err != nil {
			return err
		}

		aggregate := shared2.NewMediumAggregate(aggregateEvents)

		payload := shared2.NewMediumZurueckgegebenEvent(cmd.MediumId, cmd.NutzerId, time.Now())
		err = aggregate.HandleMediumZurueckgegeben(payload)
		if err != nil {
			return err
		}

		return v.SendEvent(payload, aggregateKey, aggregateType)
	})
}

func (v *MediumRueckgabeHandler) SendEvent(payload shared2.MediumZurueckgegebenEvent, aggregateKey, aggregateType string) error {
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
		shared2.MediumZurueckgegebenEventType,
		version+1,
		payloadJSON)

	err = v.eventStore.Save(v.ctx, event)
	if err != nil {
		return err
	}
	err = v.kafkaService.SendEvent(v.producer, shared2.MediumZurueckgegebenEventType, payloadJSON)
	if err != nil {
		panic(err)
	}
	return nil
}
