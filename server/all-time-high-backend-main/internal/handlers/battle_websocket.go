package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/GabbyWorld/all-time-high-backend/internal/logger"
	"github.com/GabbyWorld/all-time-high-backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// BattleWebSocketHandler handles battle-related WebSocket connections
type BattleWebSocketHandler struct {
	DB *gorm.DB
	// Store active connections
	Clients    map[string]*websocket.Conn
	ClientsMux sync.RWMutex
}

// NewBattleWebSocketHandler creates a new BattleWebSocketHandler
func NewBattleWebSocketHandler(db *gorm.DB) *BattleWebSocketHandler {
	return &BattleWebSocketHandler{
		DB:      db,
		Clients: make(map[string]*websocket.Conn),
	}
}

var battleUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleBattleWebSocket handles WebSocket connections for battle updates
func (h *BattleWebSocketHandler) HandleBattleWebSocket(c *gin.Context) {
	wsConn, err := battleUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Logger.Error("Battle WebSocket upgrade failed", zap.Error(err))
		return
	}
	defer wsConn.Close()

	// Register client
	h.ClientsMux.Lock()
	h.Clients[fmt.Sprintf("client_%d", time.Now().UnixNano())] = wsConn
	h.ClientsMux.Unlock()

	// Cleanup on disconnect
	defer func() {
		h.ClientsMux.Lock()
		for id, conn := range h.Clients {
			if conn == wsConn {
				delete(h.Clients, id)
				break
			}
		}
		h.ClientsMux.Unlock()
	}()

	// Set up context for cleanup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle ping/pong
	wsConn.SetPongHandler(func(string) error {
		return wsConn.SetReadDeadline(time.Now().Add(60 * time.Second))
	})

	// Start ping ticker
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := wsConn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
					logger.Logger.Warn("Failed to write ping message", zap.Error(err))
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Read messages (keep connection alive)
	for {
		_, _, err := wsConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Logger.Warn("Battle WebSocket read error", zap.Error(err))
			}
			break
		}
	}
}

// BroadcastBattleResult sends battle results to relevant clients
func (h *BattleWebSocketHandler) BroadcastBattleResult(result models.Battle) {
	h.ClientsMux.RLock()
	defer h.ClientsMux.RUnlock()

	message := struct {
		Type string        `json:"type"`
		Data models.Battle `json:"data"`
	}{
		Type: "BATTLE_RESULT",
		Data: result,
	}

	payload, err := json.Marshal(message)
	if err != nil {
		logger.Logger.Error("Failed to marshal battle result", zap.Error(err))
		return
	}

	// Send to all clients
	for _, clientConn := range h.Clients {
		err := clientConn.WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			logger.Logger.Error("Failed to send battle result to client",
				zap.Error(err))
		}
	}
}
