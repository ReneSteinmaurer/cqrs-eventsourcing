package erwerben

import (
	"context"
	"cqrs-playground/bibliothek/medien/isbn_index_projection"
	shared2 "cqrs-playground/bibliothek/medien/shared"
	"cqrs-playground/shared"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ErwerbeMediumHandler struct {
	eventStore      *shared.EventStore
	kafkaService    *shared.KafkaService
	producer        sarama.SyncProducer
	isbnIndexReader *isbn_index_projection.ISBNIndexReader
	ctx             context.Context
}

func NewErwerbeMediumHandler(
	ctx context.Context, eventStore *shared.EventStore, kafkaService *shared.KafkaService, db *pgxpool.Pool,
) *ErwerbeMediumHandler {
	producer := kafkaService.NewSyncProducer()
	isbnIndexReader := isbn_index_projection.NewPostgresISBNIndexReader(db)
	return &ErwerbeMediumHandler{
		eventStore:      eventStore,
		kafkaService:    kafkaService,
		producer:        producer,
		isbnIndexReader: isbnIndexReader,
		ctx:             ctx,
	}
}

func (n *ErwerbeMediumHandler) Handle(cmd ErwerbeMediumCommand) error {
	aggregateKey := uuid.NewString()
	aggregateType := shared2.MediumAggregateType

	if cmd.ISBN == "" {
		return errors.New("derzeit werden erwerbe ohne ISBN noch nicht unterstuetzt")
	}

	if exists, err := n.isbnIndexReader.Exists(n.ctx, cmd.ISBN); err != nil {
		return fmt.Errorf("fehler beim Prüfen der ISBN: %w", err)
	} else if exists {
		return errors.New("medium mit dieser ISBN existiert bereits")
	}

	return shared.RetryHandlerLogic(func() error {
		aggregateEvents, err := n.eventStore.GetEventsByAggregateId(n.ctx, aggregateKey, aggregateType)
		if err != nil {
			return err
		}

		payload := shared2.NewMediumErworbenEvent(cmd.ISBN, aggregateKey, cmd.Name, cmd.Genre, cmd.MediumType)

		aggregate := shared2.NewMediumAggregate(aggregateEvents)
		err = aggregate.HandleMediumErwerben(payload)
		if err != nil {
			return err
		}

		return n.SendEvent(payload, aggregateKey, aggregateType)
	})
}

func (n *ErwerbeMediumHandler) SendEvent(payload shared2.MediumErworbenEvent, aggregateKey, aggregateType string) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	version, err := n.eventStore.LoadCurrentVersion(n.ctx, aggregateKey, aggregateType)
	if err != nil {
		return err
	}

	event := shared.NewEvent(
		aggregateType,
		aggregateKey,
		shared2.MediumErworbenEventType,
		version+1,
		payloadJSON)

	err = n.eventStore.Save(n.ctx, event)
	if err != nil {
		return err
	}
	err = n.kafkaService.SendEvent(n.producer, shared2.MediumErworbenEventType, payloadJSON)
	if err != nil {
		panic(err)
	}
	return nil
}
