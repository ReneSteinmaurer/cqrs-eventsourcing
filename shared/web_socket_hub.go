package shared

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type WebSocketMessage struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type WebSocketClient struct {
	ID   string
	Conn *websocket.Conn
}

type WebSocketHub struct {
	clientsByAggregateID map[string]map[*WebSocketClient]bool
	ctx                  context.Context
	lock                 sync.RWMutex
}

func NewWebSocketHub(ctx context.Context) *WebSocketHub {
	return &WebSocketHub{
		ctx:                  ctx,
		clientsByAggregateID: make(map[string]map[*WebSocketClient]bool),
	}
}

func (hub *WebSocketHub) RegisterClientForAggregate(client *WebSocketClient, aggregateId string) {
	hub.lock.Lock()
	defer hub.lock.Unlock()

	if _, exists := hub.clientsByAggregateID[aggregateId]; !exists {
		log.Printf("[WEBSOCKET] Registering client for aggregateId=%s\n", aggregateId)
		hub.clientsByAggregateID[aggregateId] = make(map[*WebSocketClient]bool)
	}

	hub.clientsByAggregateID[aggregateId][client] = true
	log.Printf("[WEBSOCKET] Client registered for aggregateId=%s\n", aggregateId)
	go hub.listenForDisconnect(client, aggregateId)
}

func (hub *WebSocketHub) listenForDisconnect(client *WebSocketClient, aggregateId string) {
	defer client.Conn.Close()

	for {
		if _, _, err := client.Conn.ReadMessage(); err != nil {
			hub.lock.Lock()
			delete(hub.clientsByAggregateID[aggregateId], client)
			hub.lock.Unlock()
			log.Printf("[WEBSOCKET] Client disconnected for aggregateId=%s\n", aggregateId)
			break
		}
	}
}

func (hub *WebSocketHub) BroadcastToAggregate(aggregateId string, msg WebSocketMessage) error {
	hub.lock.RLock()
	defer hub.lock.RUnlock()

	clients, exists := hub.clientsByAggregateID[aggregateId]
	if !exists {
		return fmt.Errorf("[WEBSOCKET] No clients for aggregateId=%s\n", aggregateId)
	}

	for client := range clients {
		err := client.Conn.WriteJSON(msg)
		if err != nil {
			return fmt.Errorf("[WEBSOCKET] Error sending to client: %v\n", err)
		}
	}
	return nil
}

func (hub *WebSocketHub) cancel() {
	for client, _ := range hub.clientsByAggregateID {
		for cl, _ := range hub.clientsByAggregateID[client] {
			err := cl.Conn.Close()
			if err != nil {
				log.Println("[WEBSOCKET] Error closing client", err)
			}
		}
	}
}
