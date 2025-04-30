package shared

import (
	"context"
	"log"
	"time"
)

const (
	ProjectionUpdatedType = "PROJECTION_UPDATED"
)

type NotificationService struct {
	WebSocketHub *WebSocketHub
	ctx          context.Context
}

func NewNotificationService(ctx context.Context, webSocketHub *WebSocketHub) *NotificationService {
	return &NotificationService{
		WebSocketHub: webSocketHub,
		ctx:          ctx,
	}
}

func (n *NotificationService) NotifyProjectionUpdated(aggregateId string) {
	notification := WebSocketMessage{
		Type: ProjectionUpdatedType,
		Data: aggregateId,
	}

	log.Printf("[NOTIFICATION SERVICE] Notifying %s\n", aggregateId)
	Retry(5, 200*time.Millisecond, func() error {
		return n.WebSocketHub.BroadcastToAggregate(aggregateId, notification)
	})
}
