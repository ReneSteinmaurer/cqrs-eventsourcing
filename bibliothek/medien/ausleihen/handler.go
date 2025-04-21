package ausleihen

import (
	"context"
	"cqrs-playground/bibliothek/medien/bestand_projection"
	shared2 "cqrs-playground/bibliothek/medien/shared"
	"cqrs-playground/shared"
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type VerleiheMediumHandler struct {
	eventStore          *shared.EventStore
	kafkaService        *shared.KafkaService
	producer            sarama.SyncProducer
	mediumBestandReader *bestand_projection.MediumBestandReader
	leihregel           shared2.LeihregelPolicy
	ctx                 context.Context
}

func NewVerleiheMediumHandler(
	ctx context.Context, eventStore *shared.EventStore, kafkaService *shared.KafkaService, db *pgxpool.Pool,
) *VerleiheMediumHandler {
	producer := kafkaService.NewSyncProducer()
	mbr := bestand_projection.NewMediumBestandReader(db)
	return &VerleiheMediumHandler{
		eventStore:          eventStore,
		kafkaService:        kafkaService,
		producer:            producer,
		mediumBestandReader: mbr,
		leihregel:           shared2.NewStandardLeihregelPolicy(),
		ctx:                 ctx,
	}
}

func (v *VerleiheMediumHandler) Handle(cmd VerleiheMediumCommand) error {
	aggregateKey := cmd.MediumId
	aggregateType := shared2.MediumAggregateType

	if cmd.MediumId == "" || cmd.NutzerId == "" {
		return errors.New("alle Felder muessen befuellt sein")
	}

	medium, err := v.mediumBestandReader.GetByMediumId(v.ctx, cmd.MediumId)
	if err != nil {
		return err
	}

	dur := v.leihregel.DauerFuer(medium.MediumType)
	von := time.Now()
	bis := von.Add(dur)

	return shared.RetryHandlerLogic(func() error {
		aggregateEvents, err := v.eventStore.GetEventsByAggregateId(v.ctx, aggregateKey, aggregateType)
		if err != nil {
			return err
		}

		aggregate := shared2.NewMediumAggregate(aggregateEvents)

		payload := shared2.NewMediumVerliehenEvent(cmd.MediumId, cmd.NutzerId, von, bis)
		err = aggregate.HandleMediumVerleihen(payload)
		if err != nil {
			return err
		}

		return v.SendEvent(payload, aggregateKey, aggregateType)
	})
}

func (k *VerleiheMediumHandler) SendEvent(payload shared2.MediumVerliehenEvent, aggregateKey, aggregateType string) error {
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
		shared2.MediumVerliehenEventType,
		version+1,
		payloadJSON)

	err = k.eventStore.Save(k.ctx, event)
	if err != nil {
		return err
	}
	err = k.kafkaService.SendEvent(k.producer, shared2.MediumVerliehenEventType, payloadJSON)
	if err != nil {
		panic(err)
	}
	return nil
}
