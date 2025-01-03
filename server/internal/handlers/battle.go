package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/GabbyWorld/all-time-high-backend/internal/config"
	"github.com/GabbyWorld/all-time-high-backend/internal/logger"
	"github.com/GabbyWorld/all-time-high-backend/internal/models"
	"github.com/GabbyWorld/all-time-high-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BattleService struct {
	db        *gorm.DB
	wsHandler *BattleWebSocketHandler
	Config    *config.Config
}

func NewBattleService(db *gorm.DB, wsHandler *BattleWebSocketHandler, config *config.Config) *BattleService {
	return &BattleService{
		db:        db,
		wsHandler: wsHandler,
		Config:    config,
	}
}

func (s *BattleService) StartPriceMonitoring() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			logger.Logger.Info("Checking prices and triggering battles")
			s.checkPricesAndTriggerBattles()
		}
	}()
}

func (s *BattleService) checkPricesAndTriggerBattles() {
	var agents []models.Agent
	if err := s.db.Find(&agents).Error; err != nil {
		logger.Logger.Error("Failed to fetch agents", zap.Error(err))
		return
	}

	// 收集所有 TokenAddress 并批量获取相对于 SOL 的价格
	tokenAddresses := make([]string, len(agents))
	for i, agent := range agents {
		tokenAddresses[i] = agent.TokenAddress
	}

	prices, err := utils.GetMultipleTokenVsSOLPrice(tokenAddresses)
	if err != nil {
		logger.Logger.Error("Failed to get multiple token prices", zap.Error(err))
		return
	}

	for _, agent := range agents {
		price, ok := prices[agent.TokenAddress]
		if !ok {
			logger.Logger.Error("Price not found for token",
				zap.String("tokenAddress", agent.TokenAddress),
				zap.String("agentId", strconv.FormatUint(uint64(agent.ID), 10)))
			continue
		}

		// 如果 PreviousPrice 为 0，表示第一次获取价格，直接更新
		if agent.PreviousPrice == 0 {
			agent.PreviousPrice = price
			if err := s.db.Save(&agent).Error; err != nil {
				logger.Logger.Error("Failed to update agent's previous price", zap.Error(err))
			}
			continue
		}

		// 比较当前价格与 PreviousPrice
		if price > agent.PreviousPrice {
			// 当前价格高于5分钟前的价格，触发战斗
			s.triggerBattle(agent)
		}

		// 更新 PreviousPrice 为当前价格
		agent.PreviousPrice = price
		if err := s.db.Save(&agent).Error; err != nil {
			logger.Logger.Error("Failed to update agent's previous price", zap.Error(err))
		}
	}
}

func (s *BattleService) triggerBattle(attacker models.Agent) {
	// Find a random opponent
	var defender models.Agent
	if err := s.db.Where("id != ?", attacker.ID).Order("RANDOM()").First(&defender).Error; err != nil {
		logger.Logger.Error("Failed to find opponent", zap.Error(err))
		return
	}

	// Get battle outcome from ChatGPT
	battleDesc, err := utils.GenerateBattleOutcome(
		s.Config.OpenAI.APIKey,
		s.Config.OpenAI.CompletionsEndpoint,
		attacker.Name,
		attacker.Prompt,
		defender.Name,
		defender.Prompt,
	)
	if err != nil {
		logger.Logger.Error("Failed to generate battle outcome", zap.Error(err))
		return
	}

	// Determine outcome based on description
	var outcome string
	switch {
	case strings.Contains(strings.ToLower(battleDesc), "total victory"):
		outcome = "TOTAL_VICTORY"
	case strings.Contains(strings.ToLower(battleDesc), "narrow victory"):
		outcome = "NARROW_VICTORY"
	case strings.Contains(strings.ToLower(battleDesc), "crushing defeat"):
		outcome = "CRUSHING_DEFEAT"
	case strings.Contains(strings.ToLower(battleDesc), "narrow defeat"):
		outcome = "NARROW_DEFEAT"
	default:
		outcome = "NARROW_VICTORY" // Default case
	}

	// Clean up battleDesc by removing content before newlines
	if idx := strings.Index(battleDesc, "\n\n"); idx != -1 {
		battleDesc = battleDesc[idx+2:]
	} else if idx := strings.Index(battleDesc, "\n"); idx != -1 {
		battleDesc = battleDesc[idx+1:]
	}

	// Create battle result
	battle := models.Battle{
		AttackerID:  attacker.ID,
		Attacker:    attacker,
		DefenderID:  defender.ID,
		Defender:    defender,
		CreatedAt:   time.Now(),
		Outcome:     outcome,
		Description: battleDesc,
	}

	if err := s.db.Create(&battle).Error; err != nil {
		logger.Logger.Error("Failed to create battle", zap.Error(err))
		return
	}

	//update agent stats
	s.updateAgentStats(&attacker, &defender, outcome)

	if err := s.db.Preload("Attacker").Preload("Defender").First(&battle, battle.ID).Error; err != nil {
		logger.Logger.Error("Failed to retrieve created battle", zap.Error(err))
		return
	}

	// Broadcast the new result
	s.wsHandler.BroadcastBattleResult(battle)

	// Log the battle
	logger.Logger.Info("Battle completed",
		zap.String("battleId", strconv.FormatUint(uint64(battle.ID), 10)),
		zap.String("attacker", strconv.FormatUint(uint64(battle.AttackerID), 10)),
		zap.String("defender", strconv.FormatUint(uint64(battle.DefenderID), 10)),
		zap.String("outcome", battle.Outcome),
	)
}

func (s *BattleService) GetBattle(c *gin.Context) {
	var battle models.Battle
	if err := s.db.Preload("Attacker").Preload("Defender").First(&battle, c.Query("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Battle not found"})
		return
	}
	c.JSON(http.StatusOK, battle)
}

func (s *BattleService) GetBattles(c *gin.Context) {
	agentID, err := strconv.ParseUint(c.Query("agent_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid agent_id"})
		return
	}

	// get specified agent, including stats
	var agent models.Agent
	if err := s.db.First(&agent, agentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	// get all battles related to this agent
	var battles []models.Battle
	if err := s.db.Preload("Attacker").Preload("Defender").
		Where("attacker_id = ? OR defender_id = ?", agentID, agentID).
		Order("created_at desc").
		Find(&battles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch battles"})
		return
	}

	// return battle records and stats
	c.JSON(http.StatusOK, gin.H{
		"battles":  battles,
		"total":    agent.Total,
		"wins":     agent.Wins,
		"losses":   agent.Losses,
		"win_rate": agent.WinRate,
	})
}

func (s *BattleService) updateAgentStats(attacker *models.Agent, defender *models.Agent, outcome string) {
	// update total battles
	attacker.Total++
	defender.Total++

	switch outcome {
	case "TOTAL_VICTORY", "NARROW_VICTORY":
		// attacker wins
		attacker.Wins++
		defender.Losses++
	case "CRUSHING_DEFEAT", "NARROW_DEFEAT":
		// defender wins
		attacker.Losses++
		defender.Wins++
	}

	// calculate win rate
	if attacker.Total > 0 {
		attacker.WinRate = float64(attacker.Wins) / float64(attacker.Total) * 100
	}
	if defender.Total > 0 {
		defender.WinRate = float64(defender.Wins) / float64(defender.Total) * 100
	}

	// save updated agent
	if err := s.db.Save(attacker).Error; err != nil {
		logger.Logger.Error("Failed to update attacker stats", zap.Error(err))
	}
	if err := s.db.Save(defender).Error; err != nil {
		logger.Logger.Error("Failed to update defender stats", zap.Error(err))
	}
}
