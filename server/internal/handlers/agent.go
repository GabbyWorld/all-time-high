package handlers

import (
	"net/http"

	"strconv"
	"strings"
	"time"

	"github.com/GabbyWorld/all-time-high-backend/internal/config"
	"github.com/GabbyWorld/all-time-high-backend/internal/errors"
	"github.com/GabbyWorld/all-time-high-backend/internal/logger"
	"github.com/GabbyWorld/all-time-high-backend/internal/models"
	"github.com/GabbyWorld/all-time-high-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AgentHandler 处理与Agent相关的请求
type AgentHandler struct {
	DB         *gorm.DB
	Config     *config.Config
	JWTManager *utils.JWTManager
}

// AgentRequest 请求体
type AgentRequest struct {
	Name   string `json:"name" binding:"required,max=100"`
	Ticker string `json:"ticker" binding:"required,max=50"`
	Prompt string `json:"prompt" binding:"required"`
}

// AgentResponse 响应体
type AgentResponse struct {
	ID                uint      `json:"id"`
	Name              string    `json:"name"`
	Ticker            string    `json:"ticker"`
	Prompt            string    `json:"prompt"`
	Description       string    `json:"description"`
	ImageURL          string    `json:"image_url"`
	TokenAddress      string    `json:"token_address"`
	CreatedAt         time.Time `json:"created_at"`
	UserWalletAddress string    `json:"user_wallet_address"`
	MarketCap         float64   `json:"market_cap"`
	Total             int       `json:"total"`    // 新增
	Wins              int       `json:"wins"`     // 新增
	Losses            int       `json:"losses"`   // 新增
	WinRate           float64   `json:"win_rate"` // 新增
}

// CreateAgent godoc
// @Summary 创建Agent
// @Description 玩家输入name, ticker, prompt，后端生成description和图片，并保存数据。
// @Tags Agent
// @Accept  json
// @Produce  json
// @Param agent body AgentRequest true "Agent请求体"
// @Success 201 {object} AgentResponse "创建成功"
// @Failure 400 {object} errors.APIError "请求参数错误"
// @Failure 401 {object} errors.APIError "未授权"
// @Failure 500 {object} errors.APIError "服务器错误"
// @Security BearerAuth
// @Router /api/agent [post]
func (h *AgentHandler) CreateAgent(c *gin.Context) {
	var req AgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := errors.NewAPIError(errors.ErrValidation, "Request validation failed", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("CreateAgent: validation failed", zap.Error(err))
		return
	}

	// 从上下文中获取userID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		apiErr := errors.NewAPIError(errors.ErrUnauthorized, "User ID not found in context")
		c.Error(apiErr)
		logger.Logger.Error("CreateAgent: userID not found in context")
		return
	}

	userWalletAddressInterface, exists := c.Get("userWalletAddress")
	if !exists {
		apiErr := errors.NewAPIError(errors.ErrUnauthorized, "User wallet address not found in context")
		c.Error(apiErr)
		logger.Logger.Error("CreateAgent: userWalletAddress not found in context")
		return
	}
	userWalletAddress := userWalletAddressInterface.(string)

	userID, ok := userIDInterface.(uint)
	if !ok {
		apiErr := errors.NewAPIError(errors.ErrInternal, "Invalid user ID format")
		c.Error(apiErr)
		logger.Logger.Error("CreateAgent: userID format incorrect", zap.Any("userID", userIDInterface))
		return
	}

	if h.DB == nil {
		logger.Logger.Error("CreateAgent: DB is nil")
		apiErr := errors.NewAPIError(errors.ErrInternal, "Database connection not initialized")
		c.Error(apiErr)
		return
	}

	if h.Config == nil {
		logger.Logger.Error("CreateAgent: Config is nil")
		apiErr := errors.NewAPIError(errors.ErrInternal, "Configuration not loaded")
		c.Error(apiErr)
		return
	}

	// 开始数据库事务
	tx := h.DB.Begin()
	if tx.Error != nil {
		logger.Logger.Error("CreateAgent: failed to begin transaction", zap.Error(tx.Error))
		apiErr := errors.NewAPIError(errors.ErrInternal, "Failed to start transaction")
		c.Error(apiErr)
		return
	}

	// 确保事务回滚
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logger.Logger.Error("CreateAgent: panic recovered, transaction rolled back", zap.Any("panic", r))
			apiErr := errors.NewAPIError(errors.ErrInternal, "Internal server error")
			errors.RespondWithError(c, errors.GetHTTPStatus(apiErr.Code), apiErr)
		}
	}()

	// 调用ChatGPT生成描述
	description, err := utils.GenerateDescription(h.Config.OpenAI.APIKey, h.Config.OpenAI.CompletionsEndpoint, req.Name, req.Prompt)
	if err != nil {
		tx.Rollback()
		apiErr := errors.NewAPIError(errors.ErrInternal, "Failed to generate description", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("CreateAgent: failed to generate description", zap.Error(err))
		return
	}

	// 调用Stable Diffusion生成图片
	// 使用默认的 text-to-image 模式和 image/* AcceptHeader
	additionalParams := map[string]string{
		"model":         h.Config.StableDiffusion.DefaultModel,
		"output_format": h.Config.StableDiffusion.OutputFormat,
		"aspect_ratio":  h.Config.StableDiffusion.AspectRatio,
		"cfg_scale":     h.Config.StableDiffusion.Scale, // 可根据需求调整
		"seed":          "0",                            // 0表示随机种子
	}

	// 构建用于生成图像的描述模板
	imagePromptTemplate := "Pixel art of a futuristic sci-fi hero character for a 2D game, #NAME#, #PROMPT#, 32x32 pixel size, dark green matrix theme, clear top-down perspective, no angled view, isolated on transparent background, suitable for vertical scrolling shooter"

	// 用用户输入替换占位符
	imagePrompt := strings.ReplaceAll(imagePromptTemplate, "#NAME#", req.Name)
	imagePrompt = strings.ReplaceAll(imagePrompt, "#PROMPT#", req.Prompt)

	imageBytes, err := utils.GenerateImage(
		h.Config,
		imagePrompt,
		additionalParams,
	)
	if err != nil {
		tx.Rollback()
		apiErr := errors.NewAPIError(errors.ErrInternal, "Failed to generate image", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("CreateAgent: failed to generate image", zap.Error(err))
		return
	}

	// 上传图片到S3
	s3URL, err := utils.UploadImageToS3(h.Config, imageBytes)
	if err != nil {
		tx.Rollback()
		apiErr := errors.NewAPIError(errors.ErrInternal, "Failed to upload image", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("CreateAgent: failed to upload image", zap.Error(err))
		return
	}

	var tokenAddress string

	// 调用CreateToken方法（如果需要）
	if !h.Config.Solana.MockCreateToken {
		sig, err := utils.CreateToken(h.Config, s3URL, req.Name, req.Ticker, description)
		if err != nil {
			tx.Rollback()
			apiErr := errors.NewAPIError(errors.ErrTokenGeneration, "Failed to create token", err.Error())
			c.Error(apiErr)
			logger.Logger.Error("CreateAgent: failed to create token", zap.Error(err))
			return
		}

		logger.Logger.Info("CreateAgent: token created", zap.String("signature", sig.String()))

		// 等待3秒
		time.Sleep(3 * time.Second)

		// 调用GetTokenAddress方法（如果需要）
		tokenAddress, err = utils.GetTokenAddress(h.Config, sig)
		if err != nil {
			tx.Rollback()
			apiErr := errors.NewAPIError(errors.ErrInternal, "Failed to get token address", err.Error())
			c.Error(apiErr)
			logger.Logger.Error("CreateAgent: failed to get token address", zap.Error(err))
			return
		}
	} else {
		tokenAddress = "9BuQCH824VFH8XeFwXnuxjKgYW4y2joyXJYQMabspYid"
	}

	// 创建Agent记录（包括 TokenAddress）
	agent := models.Agent{
		Name:              req.Name,
		Ticker:            req.Ticker,
		Prompt:            req.Prompt,
		Description:       description,
		ImageURL:          s3URL,
		TokenAddress:      tokenAddress,
		UserID:            userID, // 关联用户ID
		CreatedAt:         time.Now(),
		HighestPrice:      2.92e-8,
		UserWalletAddress: userWalletAddress,
	}

	if err := tx.Create(&agent).Error; err != nil {
		tx.Rollback()
		apiErr := errors.NewAPIError(errors.ErrDatabase, "Failed to create agent", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("CreateAgent: failed to create agent", zap.Error(err))
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		apiErr := errors.NewAPIError(errors.ErrInternal, "Failed to commit transaction", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("CreateAgent: failed to commit transaction", zap.Error(err))
		return
	}

	// 通知 WebSocket，有新的 Agent 创建
	logger.Logger.Info("Sending agent to AgentCreatedChan", zap.Uint("agent_id", agent.ID))

	AgentCreatedChan <- agent

	marketCap := agent.HighestPrice * 1e9

	// logger.Logger.Info("CreateAgent: agent created with token address", zap.Uint("agent_id", agent.ID), zap.Uint("user_id", userID), zap.String("token_address", tokenAddress))
	c.JSON(http.StatusCreated, AgentResponse{
		ID:                agent.ID,
		Name:              agent.Name,
		Ticker:            agent.Ticker,
		Prompt:            agent.Prompt,
		Description:       agent.Description,
		ImageURL:          agent.ImageURL,
		TokenAddress:      agent.TokenAddress,
		CreatedAt:         agent.CreatedAt,
		UserWalletAddress: userWalletAddress,
		MarketCap:         marketCap,
	})
}

// AgentsResponse 响应体，包含多个Agent
type AgentsResponse struct {
	Agents []AgentResponse `json:"agents"`
}

// GetUserAgents godoc
// @Summary 获取登录用户的所有Agents（分页）
// @Description 获取当前登录用户关联的所有Agent记录，并支持分页。
// @Tags Agent
// @Produce  json
// @Param page query int false "页码(默认为1)"
// @Param page_size query int false "每页大小(默认为4)"
// @Success 200 {object} AgentsResponse "成功返回所有Agent(包含分页信息)"
// @Failure 401 {object} errors.APIError "未授权"
// @Failure 500 {object} errors.APIError "服务器错误"
// @Security BearerAuth
// @Router /api/agents [get]
// GetUserAgents 获取登录用户的所有Agents（分页）
func (h *AgentHandler) GetUserAgents(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		apiErr := errors.NewAPIError(errors.ErrUnauthorized, "User ID not found in context")
		c.Error(apiErr)
		logger.Logger.Error("GetUserAgents: userID not found in context")
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		apiErr := errors.NewAPIError(errors.ErrInternal, "Invalid user ID format")
		c.Error(apiErr)
		logger.Logger.Error("GetUserAgents: userID format incorrect", zap.Any("userID", userIDInterface))
		return
	}

	if h.DB == nil {
		logger.Logger.Error("GetUserAgents: DB is nil")
		apiErr := errors.NewAPIError(errors.ErrInternal, "Database connection not initialized")
		c.Error(apiErr)
		return
	}

	// 解析分页
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")
	page, err := utils.ParsePage(pageStr)
	if err != nil {
		page = 1
	}
	pageSize, err := utils.ParsePageSize(pageSizeStr)
	if err != nil {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// 查询总数
	var total int64
	if err := h.DB.Model(&models.Agent{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		apiErr := errors.NewAPIError(errors.ErrDatabase, "Failed to count agents", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("GetUserAgents: failed to count agents", zap.Error(err))
		return
	}

	// 如果没有查询到任何结果，直接返回
	if total == 0 {
		type AgentsPaginatedResponse struct {
			Agents   []interface{} `json:"agents"`
			Total    int64         `json:"total"`
			Page     int           `json:"page"`
			PageSize int           `json:"page_size"`
		}

		c.JSON(http.StatusOK, AgentsPaginatedResponse{
			Agents:   []interface{}{},
			Total:    0,
			Page:     page,
			PageSize: pageSize,
		})
		return
	}

	// 按创建时间降序查询
	var agents []models.Agent
	if err := h.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&agents).Error; err != nil {
		apiErr := errors.NewAPIError(errors.ErrDatabase, "Failed to retrieve agents", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("GetUserAgents: failed to retrieve agents", zap.Error(err))
		return
	}

	// 收集所有 TokenAddress 并统一获取价格
	tokenAddresses := make([]string, len(agents))
	for i, agent := range agents {
		tokenAddresses[i] = agent.TokenAddress
	}
	prices, err := utils.GetMultipleTokenPrice(tokenAddresses)
	if err != nil {
		apiErr := errors.NewAPIError(errors.ErrInternal, "Failed to get token prices", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("GetUserAgents: failed to get token prices", zap.Error(err))
		return
	}

	type agentWithFormattedTime struct {
		ID                 uint      `json:"id"`
		Name               string    `json:"name"`
		Ticker             string    `json:"ticker"`
		Prompt             string    `json:"prompt"`
		Description        string    `json:"description"`
		ImageURL           string    `json:"image_url"`
		TokenAddress       string    `json:"token_address"`
		CreatedAt          time.Time `json:"created_at"`
		MarketCap          float64   `json:"market_cap"`
		MarketCapUpdatedAt time.Time `json:"market_cap_updated_at"`
		UserWalletAddress  string    `json:"user_wallet_address"`
	}

	type AgentsPaginatedResponse struct {
		Agents   []agentWithFormattedTime `json:"agents"`
		Total    int64                    `json:"total"`
		Page     int                      `json:"page"`
		PageSize int                      `json:"page_size"`
	}

	var response AgentsPaginatedResponse

	for _, agent := range agents {
		price := prices[agent.TokenAddress]
		marketCap := price * 1e9 // 根据需求，这里假设总供应量为1e9
		response.Agents = append(response.Agents, agentWithFormattedTime{
			ID:                 agent.ID,
			Name:               agent.Name,
			Ticker:             agent.Ticker,
			Prompt:             agent.Prompt,
			Description:        agent.Description,
			ImageURL:           agent.ImageURL,
			TokenAddress:       agent.TokenAddress,
			CreatedAt:          agent.CreatedAt,
			MarketCap:          marketCap,
			MarketCapUpdatedAt: time.Now(),
			UserWalletAddress:  agent.UserWalletAddress,
		})
	}

	response.Total = total
	response.Page = page
	response.PageSize = pageSize

	c.JSON(http.StatusOK, response)
}

// GetAllAgents godoc
// @Summary 获取所有Agent（分页）
// @Description 获取数据库中所有Agent记录，并支持分页。
// @Tags Agent
// @Produce  json
// @Param page query int false "页码(默认为1)"
// @Param page_size query int false "每页大小(默认为4)"
// @Success 200 {object} AgentsResponse "成功返回所有Agent(包含分页信息)"
// @Failure 500 {object} errors.APIError "服务器错误"
// @Router /api/agents/all [get]
func (h *AgentHandler) GetAllAgents(c *gin.Context) {
	if h.DB == nil {
		logger.Logger.Error("GetAllAgents: DB is nil")
		apiErr := errors.NewAPIError(errors.ErrInternal, "Database connection not initialized")
		c.Error(apiErr)
		return
	}

	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")
	page, err := utils.ParsePage(pageStr)
	if err != nil {
		page = 1
	}
	pageSize, err := utils.ParsePageSize(pageSizeStr)
	if err != nil {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var total int64
	if err := h.DB.Model(&models.Agent{}).Count(&total).Error; err != nil {
		apiErr := errors.NewAPIError(errors.ErrDatabase, "Failed to count all agents", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("GetAllAgents: failed to count all agents", zap.Error(err))
		return
	}

	var agents []models.Agent
	// 按创建时间降序查询
	if err := h.DB.
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&agents).Error; err != nil {
		apiErr := errors.NewAPIError(errors.ErrDatabase, "Failed to retrieve all agents", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("GetAllAgents: failed to retrieve all agents", zap.Error(err))
		return
	}

	// 收集所有 TokenAddress 并统一获取价格
	tokenAddresses := make([]string, len(agents))
	for i, agent := range agents {
		tokenAddresses[i] = agent.TokenAddress
	}
	prices, err := utils.GetMultipleTokenPrice(tokenAddresses)
	if err != nil {
		apiErr := errors.NewAPIError(errors.ErrInternal, "Failed to get token prices", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("GetAllAgents: failed to get token prices", zap.Error(err))
		return
	}

	type agentWithFormattedTime struct {
		ID                 uint      `json:"id"`
		Name               string    `json:"name"`
		Ticker             string    `json:"ticker"`
		Prompt             string    `json:"prompt"`
		Description        string    `json:"description"`
		ImageURL           string    `json:"image_url"`
		TokenAddress       string    `json:"token_address"`
		CreatedAt          time.Time `json:"created_at"`
		MarketCap          float64   `json:"market_cap"`
		MarketCapUpdatedAt time.Time `json:"market_cap_updated_at"`
		UserWalletAddress  string    `json:"user_wallet_address"`
	}

	type AgentsPaginatedResponse struct {
		Agents   []agentWithFormattedTime `json:"agents"`
		Total    int64                    `json:"total"`
		Page     int                      `json:"page"`
		PageSize int                      `json:"page_size"`
	}

	var response AgentsPaginatedResponse

	for _, agent := range agents {
		price := prices[agent.TokenAddress]
		marketCap := price * 1e9 // 根据需求，这里假设总供应量为1e9

		response.Agents = append(response.Agents, agentWithFormattedTime{
			ID:                 agent.ID,
			Name:               agent.Name,
			Ticker:             agent.Ticker,
			Prompt:             agent.Prompt,
			Description:        agent.Description,
			ImageURL:           agent.ImageURL,
			TokenAddress:       agent.TokenAddress,
			CreatedAt:          agent.CreatedAt,
			MarketCap:          marketCap,
			MarketCapUpdatedAt: time.Now(),
			UserWalletAddress:  agent.UserWalletAddress,
		})
	}

	response.Total = total
	response.Page = page
	response.PageSize = pageSize

	c.JSON(http.StatusOK, response)
}

// GetAgentByID godoc
// @Summary 获取指定 ID 的 Agent
// @Description 根据 Agent ID 获取对应的 Agent 信息
// @Tags Agent
// @Produce json
// @Param id path int true "Agent ID"
// @Success 200 {object} AgentResponse "成功返回 Agent"
// @Failure 400 {object} errors.APIError "请求参数错误"
// @Failure 401 {object} errors.APIError "未授权"
// @Failure 404 {object} errors.APIError "未找到"
// @Failure 500 {object} errors.APIError "服务器错误"
// @Router /api/agent/{id} [get]
func (h *AgentHandler) GetAgentByID(c *gin.Context) {
	// 解析 URL 中的 ID 参数
	agentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apiErr := errors.NewAPIError(errors.ErrValidation, "Invalid agent ID", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("GetAgentByID: invalid agent ID", zap.Error(err))
		return
	}

	// 从数据库获取对应的 Agent
	var agent models.Agent
	if err := h.DB.First(&agent, agentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			apiErr := errors.NewAPIError(errors.ErrNotFound, "Agent not found", err.Error())
			c.Error(apiErr)
			logger.Logger.Error("GetAgentByID: agent not found", zap.Error(err))
			return
		}
		apiErr := errors.NewAPIError(errors.ErrDatabase, "Failed to retrieve agent", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("GetAgentByID: failed to retrieve agent", zap.Error(err))
		return
	}

	// 将数据库记录转换为返回体
	response := AgentResponse{
		ID:           agent.ID,
		Name:         agent.Name,
		Ticker:       agent.Ticker,
		Prompt:       agent.Prompt,
		Description:  agent.Description,
		ImageURL:     agent.ImageURL,
		TokenAddress: agent.TokenAddress,
		CreatedAt:    agent.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// 定义排行榜响应结构
type LeaderboardResponse struct {
	Leaderboard []AgentInfo `json:"leaderboard"`
}

type AgentInfo struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Ticker      string    `json:"ticker"`
	Wins        int       `json:"wins"`
	WinRate     float64   `json:"win_rate"`
	CreatedAt   time.Time `json:"created_at"`
	ImageURL    string    `json:"image_url"`
	Description string    `json:"description"`
	MarketCap   float64   `json:"market_cap"`
}

// GetLeaderboard 获取排行榜前100名的 Agent
// @Summary 获取排行榜
// @Description 获取按照胜利次数、胜率和创建时间排序的前100名 Agent
// @Tags Agent
// @Produce json
// @Success 200 {object} LeaderboardResponse "成功返回排行榜"
// @Failure 500 {object} errors.APIError "服务器错误"
// @Router /api/leaderboard [get]
func (h *AgentHandler) GetLeaderboard(c *gin.Context) {
	var agents []models.Agent

	// 查询符合条件的前100名 Agent
	if err := h.DB.
		Order("wins DESC, win_rate DESC, created_at ASC").
		Limit(100).
		Find(&agents).Error; err != nil {
		apiErr := errors.NewAPIError(errors.ErrDatabase, "Failed to get leaderboard", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("GetLeaderboard: failed to get leaderboard", zap.Error(err))
		return
	}

	// 提取所有的 TokenAddress
	tokenAddresses := make([]string, len(agents))
	for i, agent := range agents {
		tokenAddresses[i] = agent.TokenAddress
	}

	// 使用 GetMultipleTokenPrice 获取价格
	prices, err := utils.GetMultipleTokenPrice(tokenAddresses)
	if err != nil {
		apiErr := errors.NewAPIError(errors.ErrInternal, "Failed to get token prices", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("GetLeaderboard: failed to get token prices", zap.Error(err))
		return
	}

	// 构建响应数据
	var leaderboard []AgentInfo
	for _, agent := range agents {
		price := prices[agent.TokenAddress]
		marketCap := price * 1e9 // 假设总供应量为 1e9

		leaderboard = append(leaderboard, AgentInfo{
			ID:          agent.ID,
			Name:        agent.Name,
			Ticker:      agent.Ticker,
			Wins:        agent.Wins,
			WinRate:     agent.WinRate,
			CreatedAt:   agent.CreatedAt,
			ImageURL:    agent.ImageURL,
			Description: agent.Description,
			MarketCap:   marketCap,
		})
	}

	c.JSON(http.StatusOK, LeaderboardResponse{
		Leaderboard: leaderboard,
	})
}
