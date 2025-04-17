package shared

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type EventStore struct {
	DB *pgxpool.Pool
}

func (es *EventStore) Save(ctx context.Context, event Event) error {
	const query = `
		INSERT INTO events (id, type, timestamp, payload)
		VALUES ($1, $2, $3, $4)
	`

	_, err := es.DB.Exec(ctx, query, event.Id, event.Type, event.Timestamp, event.Payload)
	return err
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
