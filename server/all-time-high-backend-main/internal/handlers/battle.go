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
	ticker := time.NewTicker(15 * time.Second)
	// ticker := time.NewTicker(15 * time.Minute)
	go func() {
		for range ticker.C {
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

	for _, agent := range agents {
		price, err := utils.GetTokenPrice(agent.TokenAddress)
		if err != nil {
			logger.Logger.Error("Failed to get token price",
				zap.String("agentId", strconv.FormatUint(uint64(agent.ID), 10)),
				zap.Error(err))
			continue
		}

		if price > agent.HighestPrice {
			agent.HighestPrice = price
			if err := s.db.Save(&agent).Error; err != nil {
				logger.Logger.Error("Failed to update agent", zap.Error(err))
			}
			if agent.HighestPrice > 0 { // Only trigger battle if it's not the first price record
				s.triggerBattle(agent)
			}
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

	// Create battle result
	data := models.Battle{
		AttackerID:  attacker.ID,
		Attacker:    attacker,
		DefenderID:  defender.ID,
		Defender:    defender,
		CreatedAt:   time.Now(),
		Outcome:     outcome,
		Description: battleDesc,
	}

	if err := s.db.Create(&data).Error; err != nil {
		logger.Logger.Error("Failed to create battle", zap.Error(err))
		return
	}

	// Broadcast the new result
	s.wsHandler.BroadcastBattleResult(data)

	// Log the battle
	logger.Logger.Info("Battle completed",
		zap.String("battleId", strconv.FormatUint(uint64(data.ID), 10)),
		zap.String("attacker", strconv.FormatUint(uint64(data.AttackerID), 10)),
		zap.String("defender", strconv.FormatUint(uint64(data.DefenderID), 10)),
		zap.String("outcome", data.Outcome),
	)
}

func (s *BattleService) GetBattle(c *gin.Context) {
	var battle models.Battle
	if err := s.db.First(&battle, c.Query("id")).Error; err != nil {
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
	var battles []models.Battle
	if err := s.db.Where("attacker_id = ? OR defender_id = ?", agentID, agentID).Order("created_at desc").Find(&battles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch battles"})
		return
	}

	// Calculate statistics
	totalBattles := len(battles)
	winCount := 0
	loseCount := 0

	for _, battle := range battles {
		if battle.Outcome == "TOTAL_VICTORY" || battle.Outcome == "NARROW_VICTORY" {
			if battle.AttackerID == uint(agentID) {
				winCount++
			} else {
				loseCount++
			}
		} else {
			if battle.AttackerID == uint(agentID) {
				loseCount++
			} else {
				winCount++
			}
		}
	}

	winRate := 0.0
	if totalBattles > 0 {
		winRate = float64(winCount) / float64(totalBattles) * 100
	}

	// Return battles and statistics
	c.JSON(http.StatusOK, gin.H{
		"battles":  battles,
		"total":    totalBattles,
		"wins":     winCount,
		"losses":   loseCount,
		"win_rate": winRate,
	})
}
