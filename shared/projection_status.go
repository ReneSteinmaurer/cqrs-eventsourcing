package shared

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"sync"
	"time"
)

type ProjectionStatus struct {
	ProjectionName         string    `DB:"projection_name"`
	LastProcessedEventId   string    `DB:"last_processed_event_id"`
	LastProcessedTimestamp time.Time `DB:"last_processed_timestamp"`
}

type ProjectionStateUpdater struct {
	DB    *pgxpool.Pool
	ctx   context.Context
	mutex sync.Mutex
}

func NewProjectionStateUpdater(ctx context.Context, db *pgxpool.Pool) *ProjectionStateUpdater {
	return &ProjectionStateUpdater{
		DB:    db,
		ctx:   ctx,
		mutex: sync.Mutex{},
	}
}

func (cp *ProjectionStateUpdater) UpdateProjectionState(projectionName, eventId string, time time.Time) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	const query = `
		  insert into projection_status (projection_name, last_processed_event_id, last_processed_timestamp) 
			values ($1, $2, $3)
			on conflict (projection_name) 
			do update set 
			last_processed_event_id = $2,
			last_processed_timestamp = $3
	`
	_, err := cp.DB.Exec(cp.ctx, query, projectionName, eventId, time)
	if err != nil {
		log.Println("Error updating projection_status table", err)
	}
}
