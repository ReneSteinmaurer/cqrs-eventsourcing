package shared

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type EventStore struct {
	DB *pgxpool.Pool
}

func (es *EventStore) Save(ctx context.Context, event Event) error {
	const query = `
		INSERT INTO events (id, type, timestamp, payload, aggregate_key, aggregate_type, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := es.DB.Exec(
		ctx,
		query,
		event.Id,
		event.Type,
		event.Timestamp,
		event.Payload,
		event.AggregateKey,
		event.AggregateType,
		event.Version)

	if IsVersionConflict(err) {
		fmt.Printf("Postgres error: event could not be persisted due to version mismatch aggregateKey=%s aggregateType=%s version=%d\n",
			event.AggregateKey, event.AggregateType, event.Version)
	}
	return err
}

func (es *EventStore) LoadCurrentVersion(ctx context.Context, aggregateKey, aggregateType string) (int, error) {
	const query = `SELECT COALESCE(MAX(version), 0) FROM events WHERE aggregate_key = $1 and aggregate_type = $2`
	var version int
	err := es.DB.QueryRow(ctx, query, aggregateKey, aggregateType).Scan(&version)
	return version, err
}

func (es *EventStore) GetEventsByAggregateId(ctx context.Context, aggregateKey, aggregateType string) ([]Event, error) {
	const query = `
		SELECT id, type, timestamp, payload, aggregate_key, aggregate_type, version
		FROM events
		WHERE aggregate_key = $1
		AND aggregate_type = $2
		ORDER BY version ASC
	`

	rows, err := es.DB.Query(ctx, query, aggregateKey, aggregateType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.Id, &e.Type, &e.Timestamp, &e.Payload, &e.AggregateType, &e.AggregateKey, &e.Version); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (es *EventStore) GetEventsSince(ctx context.Context, lastUpdatedTime time.Time) ([]Event, error) {
	const query = `
		SELECT id, type, timestamp, payload
		FROM events
		WHERE timestamp > $1
		ORDER BY timestamp ASC
	`

	rows, err := es.DB.Query(ctx, query, lastUpdatedTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		if err := rows.Scan(&event.Id, &event.Type, &event.Timestamp, &event.Payload); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}

func (es *EventStore) GetLastUpdateFromProjection(ctx context.Context, projectionName string) (ProjectionStatus, error) {
	const query = `
		select projection_name, last_processed_event_id, last_processed_timestamp
		from projection_status
		where projection_name = $1
	`

	var projection ProjectionStatus
	err := es.DB.QueryRow(ctx, query, projectionName).Scan(
		&projection.ProjectionName,
		&projection.LastProcessedEventId,
		&projection.LastProcessedTimestamp,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ProjectionStatus{
				ProjectionName:         projectionName,
				LastProcessedEventId:   "",
				LastProcessedTimestamp: time.Unix(0, 0),
			}, nil
		}
		return ProjectionStatus{}, err
	}

	return projection, nil
}
