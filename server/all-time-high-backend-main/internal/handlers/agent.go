package handlers

import (
	"net/http"

	"github.com/GabbyWorld/all-time-high-backend/internal/config"
	"github.com/GabbyWorld/all-time-high-backend/internal/errors"
	"github.com/GabbyWorld/all-time-high-backend/internal/logger"
	"github.com/GabbyWorld/all-time-high-backend/internal/models"
	"github.com/GabbyWorld/all-time-high-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
	"time"
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
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Ticker       string    `json:"ticker"`
	Prompt       string    `json:"prompt"`
	Description  string    `json:"description"`
	ImageURL     string    `json:"image_url"`
	TokenAddress string    `json:"token_address"` // 新增字段
	CreatedAt    time.Time `json:"created_at"`
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
		errors.RespondWithError(c, errors.GetHTTPStatus(apiErr.Code), apiErr)
		return
	}

	// 从上下文中获取userID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		apiErr := errors.NewAPIError(errors.ErrUnauthorized, "User ID not found in context")
		c.Error(apiErr)
		logger.Logger.Error("CreateAgent: userID not found in context")
		errors.RespondWithError(c, errors.GetHTTPStatus(apiErr.Code), apiErr)
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		apiErr := errors.NewAPIError(errors.ErrInternal, "Invalid user ID format")
		c.Error(apiErr)
		logger.Logger.Error("CreateAgent: userID format incorrect", zap.Any("userID", userIDInterface))
		errors.RespondWithError(c, errors.GetHTTPStatus(apiErr.Code), apiErr)
		return
	}

	if h.DB == nil {
		logger.Logger.Error("CreateAgent: DB is nil")
		apiErr := errors.NewAPIError(errors.ErrInternal, "Database connection not initialized")
		c.Error(apiErr)
		errors.RespondWithError(c, errors.GetHTTPStatus(apiErr.Code), apiErr)
		return
	}

	if h.Config == nil {
		logger.Logger.Error("CreateAgent: Config is nil")
		apiErr := errors.NewAPIError(errors.ErrInternal, "Configuration not loaded")
		c.Error(apiErr)
		errors.RespondWithError(c, errors.GetHTTPStatus(apiErr.Code), apiErr)
		return
	}

	// 开始数据库事务
	tx := h.DB.Begin()
	if tx.Error != nil {
		logger.Logger.Error("CreateAgent: failed to begin transaction", zap.Error(tx.Error))
		apiErr := errors.NewAPIError(errors.ErrInternal, "Failed to start transaction")
		c.Error(apiErr)
		errors.RespondWithError(c, errors.GetHTTPStatus(apiErr.Code), apiErr)
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
		errors.RespondWithError(c, errors.GetHTTPStatus(apiErr.Code), apiErr)
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
		errors.RespondWithError(c, errors.GetHTTPStatus(apiErr.Code), apiErr)
		return
	}

	// 上传图片到S3
	s3URL, err := utils.UploadImageToS3(h.Config, imageBytes)
	if err != nil {
		tx.Rollback()
		apiErr := errors.NewAPIError(errors.ErrInternal, "Failed to upload image", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("CreateAgent: failed to upload image", zap.Error(err))
		errors.RespondWithError(c, errors.GetHTTPStatus(apiErr.Code), apiErr)
		return
	}

	// 调用CreateToken方法（如果需要）
	// sig, err := utils.CreateToken(h.Config, s3URL)
	// if err != nil {
	// 	tx.Rollback()
	// 	apiErr := errors.NewAPIError(errors.ErrTokenGeneration, "Failed to create token", err.Error())
	// 	c.Error(apiErr)
	// 	logger.Logger.Error("CreateAgent: failed to create token", zap.Error(err))
	// 	errors.RespondWithError(c, errors.GetHTTPStatus(apiErr.Code), apiErr)
	// 	return
	// }

	// 调用GetTokenAddress方法（如果需要）
	// tokenAddress, err := utils.GetTokenAddress(h.Config, sig)
	// if err != nil {
	// 	tx.Rollback()
	// 	apiErr := errors.NewAPIError(errors.ErrInternal, "Failed to get token address", err.Error())
	// 	c.Error(apiErr)
	// 	logger.Logger.Error("CreateAgent: failed to get token address", zap.Error(err))
	// 	errors.RespondWithError(c, errors.GetHTTPStatus(apiErr.Code), apiErr)
	// 	return
	// }

	// 创建Agent记录（包括 TokenAddress）
	agent := models.Agent{
		Name:         req.Name,
		Ticker:       req.Ticker,
		Prompt:       req.Prompt,
		Description:  description,
		ImageURL:     s3URL,
		TokenAddress: "4cYAZQsLkjTeqLsJr5miYXwnMFr59QgUe6Gkp8sgSmB3", // 直接赋值
		UserID:       userID,                                         // 关联用户ID
		CreatedAt:    time.Now(),
	}

	if err := tx.Create(&agent).Error; err != nil {
		tx.Rollback()
		apiErr := errors.NewAPIError(errors.ErrDatabase, "Failed to create agent", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("CreateAgent: failed to create agent", zap.Error(err))
		errors.RespondWithError(c, errors.GetHTTPStatus(apiErr.Code), apiErr)
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		apiErr := errors.NewAPIError(errors.ErrInternal, "Failed to commit transaction", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("CreateAgent: failed to commit transaction", zap.Error(err))
		errors.RespondWithError(c, errors.GetHTTPStatus(apiErr.Code), apiErr)
		return
	}

	// 通知 WebSocket，有新的 Agent 创建
	logger.Logger.Info("Sending agent to AgentCreatedChan", zap.Uint("agent_id", agent.ID))

	AgentCreatedChan <- agent

	// logger.Logger.Info("CreateAgent: agent created with token address", zap.Uint("agent_id", agent.ID), zap.Uint("user_id", userID), zap.String("token_address", tokenAddress))
	c.JSON(http.StatusCreated, AgentResponse{
		ID:           agent.ID,
		Name:         agent.Name,
		Ticker:       agent.Ticker,
		Prompt:       agent.Prompt,
		Description:  agent.Description,
		ImageURL:     agent.ImageURL,
		TokenAddress: agent.TokenAddress,
		CreatedAt:    agent.CreatedAt,
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
		errors.RespondWithError(c, errors.GetHTTPStatus(apiErr.Code), apiErr)
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		apiErr := errors.NewAPIError(errors.ErrInternal, "Invalid user ID format")
		c.Error(apiErr)
		logger.Logger.Error("GetUserAgents: userID format incorrect", zap.Any("userID", userIDInterface))
		errors.RespondWithError(c, errors.GetHTTPStatus(apiErr.Code), apiErr)
		return
	}

	if h.DB == nil {
		logger.Logger.Error("GetUserAgents: DB is nil")
		apiErr := errors.NewAPIError(errors.ErrInternal, "Database connection not initialized")
		c.Error(apiErr)
		errors.RespondWithError(c, errors.GetHTTPStatus(apiErr.Code), apiErr)
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
		errors.RespondWithError(c, errors.GetHTTPStatus(apiErr.Code), apiErr)
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
		errors.RespondWithError(c, errors.GetHTTPStatus(apiErr.Code), apiErr)
		return
	}

	// 用于返回的结构体
	type agentWithFormattedTime struct {
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
	}

	type AgentsPaginatedResponse struct {
		Agents   []agentWithFormattedTime `json:"agents"`
		Total    int64                    `json:"total"`
		Page     int                      `json:"page"`
		PageSize int                      `json:"page_size"`
	}

	var response AgentsPaginatedResponse

	for _, agent := range agents {
		// 获取最新价格
		price, _ := utils.GetTokenPrice(agent.TokenAddress) // 根据需要处理错误

		// MarketCap = price * 10^9
		marketCap := price * 1e9

		response.Agents = append(response.Agents, agentWithFormattedTime{
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
// @Security BearerAuth
// @Router /api/agents/all [get]
func (h *AgentHandler) GetAllAgents(c *gin.Context) {
	if h.DB == nil {
		logger.Logger.Error("GetAllAgents: DB is nil")
		apiErr := errors.NewAPIError(errors.ErrInternal, "Database connection not initialized")
		c.Error(apiErr)
		errors.RespondWithError(c, errors.GetHTTPStatus(apiErr.Code), apiErr)
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
		errors.RespondWithError(c, errors.GetHTTPStatus(apiErr.Code), apiErr)
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
		errors.RespondWithError(c, errors.GetHTTPStatus(apiErr.Code), apiErr)
		return
	}

	// 与 GetUserAgents 类似，增加 MarketCap、MarketCapUpdatedAt
	type agentWithFormattedTime struct {
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
	}

	type AgentsPaginatedResponse struct {
		Agents   []agentWithFormattedTime `json:"agents"`
		Total    int64                    `json:"total"`
		Page     int                      `json:"page"`
		PageSize int                      `json:"page_size"`
	}

	var response AgentsPaginatedResponse

	for _, agent := range agents {
		price, _ := utils.GetTokenPrice(agent.TokenAddress)
		marketCap := price * 1e9
		response.Agents = append(response.Agents, agentWithFormattedTime{
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
		})
	}

	response.Total = total
	response.Page = page
	response.PageSize = pageSize

	c.JSON(http.StatusOK, response)
}
