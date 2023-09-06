package handlers

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/haydnmeyburgh/cod-eSports-tournament-manager/internal/realtimemanager"
)

type WebSocketHandler struct {
	WebSocketHub *realtimemanager.WebSocketHub
}

func (wh *WebSocketHandler) HandleWebSocketMessages(c *gin.Context, conn *websocket.Conn) {
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		var message map[string]interface{}
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Printf("WebSocket message unmarshal error: %v", err)
			continue
		}

		action, ok := message["action"].(string)
		if !ok {
			log.Printf("Invalid WebSocket message format")
			continue
		}

		// Handle WebSocket messages based on the "action"
		switch action {
		case "tournament_created":
			// Handle tournament creation and broadcast to clients
			wh.handleTournamentCreated(msg)
		case "tournament_updated":
			// Handle tournament updates and broadcast to clients
			wh.handleTournamentUpdated(msg)
		case "match_result_created":
			// Handle match result creation and broadcast to clients
			wh.handleMatchResultCreated(msg)
		case "match_result_updated":
			// Handle match result updates and broadcast to clients
			wh.handleMatchResultUpdated(msg)
		case "team_created":
			// Handle team creation and broadcast to clients
			wh.handleTeamCreated(msg)
		case "team_updated":
			// Handle team updates and broadcast to clients
			wh.handleTeamUpdated(msg)
		case "match_created":
			// Handle match creation and broadcast to clients
			wh.handleMatchCreated(msg)
		case "match_updated":
			// Handle match updates and broadcast to clients
			wh.handleMatchUpdated(msg)
		default:
			log.Printf("Unknown WebSocket action: %s", action)
		}
	}
}

// Define handler functions for each message type
func (wh *WebSocketHandler) handleTournamentCreated(msg []byte) {
	// Handle tournament creation and broadcast to clients
	wh.WebSocketHub.Broadcast(msg)
}

func (wh *WebSocketHandler) handleTournamentUpdated(msg []byte) {
	// Handle tournament updates and broadcast to clients
	wh.WebSocketHub.Broadcast(msg)
}

func (wh *WebSocketHandler) handleMatchResultCreated(msg []byte) {
	// Handle match result creation and broadcast to clients
	wh.WebSocketHub.Broadcast(msg)
}

func (wh *WebSocketHandler) handleMatchResultUpdated(msg []byte) {
	// Handle match result updates and broadcast to clients
	wh.WebSocketHub.Broadcast(msg)
}

func (wh *WebSocketHandler) handleTeamCreated(msg []byte) {
	// Handle team creation and broadcast to clients
	wh.WebSocketHub.Broadcast(msg)
}

func (wh *WebSocketHandler) handleTeamUpdated(msg []byte) {
	// Handle team updates and broadcast to clients
	wh.WebSocketHub.Broadcast(msg)
}

func (wh *WebSocketHandler) handleMatchCreated(msg []byte) {
	// Handle match creation and broadcast to clients
	wh.WebSocketHub.Broadcast(msg)
}

func (wh *WebSocketHandler) handleMatchUpdated(msg []byte) {
	// Handle match updates and broadcast to clients
	wh.WebSocketHub.Broadcast(msg)
}

func NewWebSocketHandler(webSocketHub *realtimemanager.WebSocketHub) *WebSocketHandler {
	return &WebSocketHandler{
		WebSocketHub: webSocketHub,
	}
}
