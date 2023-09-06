package realtimemanager

import (
	"log"
	"sync"
	"github.com/gorilla/websocket"
)

// Manages WebSocket clients
type WebSocketHub struct {
	clients map[*websocket.Conn]struct{}
	mu sync.Mutex
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		clients: make(map[*websocket.Conn]struct{}),
	}
}

// Add client
func (wh *WebSocketHub) AddClient(client *websocket.Conn) {
	wh.mu.Lock()
	defer wh.mu.Unlock()
	wh.clients[client] = struct{}{}
}

// Remove client
func (wh *WebSocketHub) RemoveClient(client *websocket.Conn) {
	wh.mu.Lock()
	defer wh.mu.Unlock()
	delete(wh.clients, client)
}

// Broadcast a message
func (wh *WebSocketHub) Broadcast(message []byte) {
	wh.mu.Lock()
	defer wh.mu.Unlock()
	for client := range wh.clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Error sending WebSocket message: %v", err)
			client.Close()
			delete(wh.clients, client)
		}
	}
}