package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/GabbyWorld/all-time-high-backend/internal/logger"
	"github.com/GabbyWorld/all-time-high-backend/internal/models"
	"github.com/GabbyWorld/all-time-high-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ======================================================
// ========== 1. 全局事件源 & 广播器定义与启动 ============
// ======================================================

// 全局或单例的 AgentCreatedChan（只在一个地方定义）
var AgentCreatedChan = make(chan models.Agent, 10)

// AgentBroadcaster 用于管理所有订阅者，并把 Agent 事件广播给他们
type AgentBroadcaster struct {
	mu          sync.RWMutex
	subscribers map[chan models.Agent]struct{}
}

var GlobalBroadcaster = &AgentBroadcaster{
	subscribers: make(map[chan models.Agent]struct{}),
}

// Start 在后台启动一个 goroutine，
// 统一从 AgentCreatedChan 读取新的 Agent，然后广播给所有订阅者
func (b *AgentBroadcaster) Start() {
	go func() {
		for agent := range AgentCreatedChan {
			b.mu.RLock()
			for sub := range b.subscribers {
				// 非阻塞写，防止某个订阅者不读时阻塞整个广播
				select {
				case sub <- agent:
				default:
				}
			}
			b.mu.RUnlock()
		}
	}()
}

// Subscribe 用于注册一个新的订阅 channel，并返回它
func (b *AgentBroadcaster) Subscribe() chan models.Agent {
	ch := make(chan models.Agent, 10)
	b.mu.Lock()
	b.subscribers[ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

// Unsubscribe 用于取消订阅并关闭 channel
func (b *AgentBroadcaster) Unsubscribe(ch chan models.Agent) {
	b.mu.Lock()
	delete(b.subscribers, ch)
	close(ch)
	b.mu.Unlock()
}

// 你可以在 init() 或者 main() 函数中调用，保证程序启动时就开启广播逻辑
func init() {
	GlobalBroadcaster.Start()
}

// ======================================================
// =========== 2. WebSocket Handler 逻辑示例 ============
// ======================================================

type AgentWebSocketHandler struct {
	DB *gorm.DB
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleAgentWebSocket 用于处理新的 WebSocket 连接
func (h *AgentWebSocketHandler) HandleAgentWebSocket(c *gin.Context) {
	// 升级 HTTP 为 WebSocket
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Logger.Error("WebSocket upgrade failed", zap.Error(err))
		return
	}
	defer wsConn.Close()

	// 【订阅】获取一个专属的 channel，用于接收 Agent 广播
	agentChan := GlobalBroadcaster.Subscribe()
	// 在函数退出（连接结束）时，取消订阅，防止内存泄漏
	defer GlobalBroadcaster.Unsubscribe(agentChan)

	// 设置上下文，用于控制 goroutine 生命周期
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建一个互斥锁来同步 WebSocket 写操作
	var writeMutex sync.Mutex

	// 设置读超时
	wsConn.SetReadDeadline(time.Now().Add(2 * time.Minute))

	// 当服务端收到 Pong 消息时，延长读超时
	wsConn.SetPongHandler(func(appData string) error {
		logger.Logger.Debug("Received Pong from client")
		wsConn.SetReadDeadline(time.Now().Add(2 * time.Minute))
		return nil
	})

	// 当服务端收到 Ping 消息时，延长读超时并发送 Pong
	wsConn.SetPingHandler(func(appData string) error {
		logger.Logger.Debug("Received Ping from client")
		wsConn.SetReadDeadline(time.Now().Add(2 * time.Minute))
		writeMutex.Lock()
		defer writeMutex.Unlock()
		return wsConn.WriteMessage(websocket.PongMessage, []byte(appData))
	})

	// 定期发送 Ping 消息给客户端，保持连接活跃
	pingTicker := time.NewTicker(30 * time.Second)
	defer pingTicker.Stop()

	// goroutine：负责定期 Ping 客户端
	go func() {
		for {
			select {
			case <-pingTicker.C:
				writeMutex.Lock()
				err := wsConn.WriteMessage(websocket.PingMessage, nil)
				writeMutex.Unlock()
				if err != nil {
					logger.Logger.Warn("Failed to send Ping", zap.Error(err))
					cancel() // 出现错误时，取消上下文
					return
				}
				logger.Logger.Debug("Sent Ping to client")
			case <-ctx.Done():
				return
			}
		}
	}()

	// goroutine：不断读取 agentChan，并把新的 Agent 发送给客户端
	go func() {
		for {
			select {
			case agent, ok := <-agentChan:
				if !ok {
					// 说明 agentChan 已被关闭或取消订阅
					return
				}

				// 从第三方函数获取当前 Token 价格
				price, err := utils.GetTokenPrice(agent.TokenAddress)
				if err != nil {
					logger.Logger.Error("Failed to get token price", zap.Error(err))
					continue
				}
				marketCap := price * 1e9

				// 这里对 Agent 做一次简单格式化，以便发送到前端
				formattedAgent := struct {
					ID                 uint    `json:"id"`
					Name               string  `json:"name"`
					Ticker             string  `json:"ticker"`
					Prompt             string  `json:"prompt"`
					Description        string  `json:"description"`
					ImageURL           string  `json:"image_url"`
					TokenAddress       string  `json:"token_address"`
					CreatedAt          string  `json:"created_at"`
					MarketCap          float64 `json:"market_cap"`
					MarketCapUpdatedAt string  `json:"market_cap_updated_at"`
				}{
					ID:                 agent.ID,
					Name:               agent.Name,
					Ticker:             agent.Ticker,
					Prompt:             agent.Prompt,
					Description:        agent.Description,
					ImageURL:           agent.ImageURL,
					TokenAddress:       agent.TokenAddress,
					CreatedAt:          agent.CreatedAt.Format("2006-01-02 15:04:05"),
					MarketCap:          marketCap,
					MarketCapUpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
				}

				// 序列化成 JSON
				data, marshalErr := json.Marshal(formattedAgent)
				if marshalErr != nil {
					logger.Logger.Error("JSON marshal error for new Agent", zap.Error(marshalErr))
					continue
				}

				logger.Logger.Info("Sending agent to client", zap.Uint("agent_id", agent.ID))

				// 通过 WebSocket 发送给前端
				writeMutex.Lock()
				writeErr := wsConn.WriteMessage(websocket.TextMessage, data)
				writeMutex.Unlock()
				if writeErr != nil {
					logger.Logger.Warn("Failed to write message to WebSocket", zap.Error(writeErr))
					cancel() // 关闭上下文
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// goroutine：侦听客户端消息，若客户端主动关闭，退出循环
	go func() {
		defer cancel()
		for {
			// 不断读取消息(仅用于保持连接，也可根据需要处理客户端数据)
			if _, _, err := wsConn.ReadMessage(); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logger.Logger.Warn("WebSocket read error", zap.Error(err))
				} else {
					logger.Logger.Info("WebSocket connection closed by client", zap.Error(err))
				}
				return
			}
		}
	}()

	// 阻塞等待上下文结束（或者你也可以在此使用 select{}）
	<-ctx.Done()

	// ---------------------------
	// 清理资源：关闭连接、发送关闭消息
	// ---------------------------
	writeMutex.Lock()
	closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Connection closed")
	_ = wsConn.WriteMessage(websocket.CloseMessage, closeMsg)
	writeMutex.Unlock()

	logger.Logger.Info("WebSocket connection closed gracefully")
}
