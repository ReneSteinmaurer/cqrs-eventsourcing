package ws

import (
	"context"
	"cqrs-playground/shared"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type NotificationWsAPI struct {
	websocketHub *shared.WebSocketHub
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewNotificationWsAPI(ctx context.Context, websocketHub *shared.WebSocketHub) *NotificationWsAPI {
	ctx, cancel := context.WithCancel(ctx)
	return &NotificationWsAPI{
		websocketHub: websocketHub,
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (n *NotificationWsAPI) Handle(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	aggregateId := params.Get("aggregateId")
	log.Println("Websocket connection for aggregateId: ", aggregateId)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading:", err)
		return
	}

	client := &shared.WebSocketClient{
		ID:   aggregateId,
		Conn: conn,
	}

	n.websocketHub.RegisterClientForAggregate(client, aggregateId)
}
