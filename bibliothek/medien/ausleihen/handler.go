package ausleihen

import (
	"context"
	"cqrs-playground/bibliothek/medien/ausleihen/events"
	"cqrs-playground/bibliothek/medien/projections/bestand"
	shared2 "cqrs-playground/bibliothek/medien/shared"
	"cqrs-playground/bibliothek/nutzer/projections/nutzer"
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
	mediumBestandReader *bestand.MediumBestandReader
	nutzerReader        *nutzer.NutzerReader
	leihregel           LeihregelPolicy
	ctx                 context.Context
}

func NewVerleiheMediumHandler(
	ctx context.Context, eventStore *shared.EventStore, kafkaService *shared.KafkaService, db *pgxpool.Pool,
) *VerleiheMediumHandler {
	producer := kafkaService.NewSyncProducer()
	mbr := bestand.NewMediumBestandReader(db)
	nr := nutzer.NewNutzerReader(db)
	return &VerleiheMediumHandler{
		eventStore:          eventStore,
		kafkaService:        kafkaService,
		producer:            producer,
		mediumBestandReader: mbr,
		nutzerReader:        nr,
		leihregel:           NewStandardLeihregelPolicy(),
		ctx:                 ctx,
	}
}

func (v *VerleiheMediumHandler) Handle(cmd VerleiheMediumCommand) (string, error) {
	aggregateKey := cmd.MediumId
	aggregateType := shared2.MediumAggregateType

	if cmd.MediumId == "" || cmd.NutzerId == "" {
		return "", errors.New("alle Felder muessen befuellt sein")
	}

	medium, err := v.mediumBestandReader.GetByMediumId(v.ctx, cmd.MediumId)
	if err != nil {
		return "", err
	}

	nutzerExists, err := v.nutzerReader.Exists(v.ctx, cmd.NutzerId)
	if err != nil {
		return "", err
	}

	if !nutzerExists {
		return "", errors.New("nutzer existiert nicht")
	}

	dur := v.leihregel.DauerFuer(medium.MediumType)
	von := time.Now()
	bis := von.Add(dur)

	return aggregateKey, shared.RetryHandlerBasedOnVersionConflict(func() error {
		aggregateEvents, err := v.eventStore.GetEventsByAggregateId(v.ctx, aggregateKey, aggregateType)
		if err != nil {
			return err
		}

		aggregate := shared2.NewMediumAggregate(aggregateEvents)

		payload := events.NewMediumVerliehenEvent(cmd.MediumId, cmd.NutzerId, von, bis)
		err = aggregate.HandleMediumVerleihen(payload)
		if err != nil {
			return err
		}

		return v.SendEvent(payload, aggregateKey, aggregateType)
	})
}

func (v *VerleiheMediumHandler) SendEvent(payload events.MediumVerliehenEvent, aggregateKey, aggregateType string) error {
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
		events.MediumVerliehenEventType,
		version+1,
		payloadJSON)

	err = v.eventStore.Save(v.ctx, event)
	if err != nil {
		return err
	}
	err = v.kafkaService.SendEvent(v.producer, events.MediumVerliehenEventType, payloadJSON)
	if err != nil {
		panic(err)
	}
	return nil
}
