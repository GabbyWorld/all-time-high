// internal/handlers/user.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/GabbyWorld/all-time-high-backend/internal/errors"
	"github.com/GabbyWorld/all-time-high-backend/internal/logger"
	"github.com/GabbyWorld/all-time-high-backend/internal/models"
	"github.com/GabbyWorld/all-time-high-backend/pkg/utils"
)

type UserHandler struct {
	DB         *gorm.DB
	JWTManager *utils.JWTManager
}

// ConnectWalletRequest 用户连接Phantom钱包的请求体
type ConnectWalletRequest struct {
	// WalletAddress 用户的Phantom钱包地址
	WalletAddress string `json:"wallet_address" binding:"required"`
	// Username 用户名（可选）
	Username string `json:"username" binding:"omitempty,max=50"`
}

// ConnectWalletResponse 用户连接钱包的响应体
type ConnectWalletResponse struct {
	// Message 操作消息
	Message string `json:"message"`
	// User 用户信息
	User models.User `json:"user"`
	// Token JWT令牌
	Token string `json:"token"`
}

// ConnectWallet godoc
// @Summary 连接Phantom钱包
// @Description 用户通过Phantom钱包连接到系统，创建或更新用户信息，并返回JWT令牌。
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param connect_wallet body ConnectWalletRequest true "用户连接Phantom钱包请求体"
// @Success 201 {object} ConnectWalletResponse "用户已创建"
// @Success 200 {object} ConnectWalletResponse "用户已连接"
// @Failure 400 {object} errors.APIError "请求参数错误"
// @Failure 500 {object} errors.APIError "服务器错误"
// @Router /api/connect_wallet [post]
func (h *UserHandler) ConnectWallet(c *gin.Context) {
	var req ConnectWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := errors.NewAPIError(errors.ErrValidation, "Request validation failed", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("ConnectWallet: validation failed", zap.Error(err))
		return
	}

	var user models.User
	// 查找是否已有该钱包地址的用户
	result := h.DB.Where("wallet_address = ?", req.WalletAddress).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// 创建新用户
			user = models.User{
				WalletAddress: req.WalletAddress,
				Username:      req.Username,
			}
			if err := h.DB.Create(&user).Error; err != nil {
				apiErr := errors.NewAPIError(errors.ErrDatabase, "Failed to create user", err.Error())
				c.Error(apiErr)
				logger.Logger.Error("ConnectWallet: failed to create user", zap.Error(err))
				return
			}
			// 生成JWT
			token, err := h.JWTManager.GenerateToken(&user)
			if err != nil {
				apiErr := errors.NewAPIError(errors.ErrTokenGeneration, "Failed to generate token", err.Error())
				c.Error(apiErr)
				logger.Logger.Error("ConnectWallet: failed to generate token", zap.Error(err))
				return
			}
			logger.Logger.Info("ConnectWallet: user created", zap.Uint("user_id", user.ID), zap.String("wallet_address", user.WalletAddress))
			c.JSON(http.StatusCreated, ConnectWalletResponse{
				Message: "User created",
				User:    user,
				Token:   token,
			})
			return
		}
		// 其他错误
		apiErr := errors.NewAPIError(errors.ErrDatabase, "Database query error", result.Error.Error())
		c.Error(apiErr)
		logger.Logger.Error("ConnectWallet: database query error", zap.Error(result.Error))
		return
	}

	// 更新用户名（如果提供）
	if req.Username != "" && req.Username != user.Username {
		user.Username = req.Username
		if err := h.DB.Save(&user).Error; err != nil {
			apiErr := errors.NewAPIError(errors.ErrDatabase, "Failed to update user", err.Error())
			c.Error(apiErr)
			logger.Logger.Error("ConnectWallet: failed to update user", zap.Error(err))
			return
		}
		logger.Logger.Info("ConnectWallet: user updated", zap.Uint("user_id", user.ID), zap.String("username", user.Username))
	}

	// 生成JWT
	token, err := h.JWTManager.GenerateToken(&user)
	if err != nil {
		apiErr := errors.NewAPIError(errors.ErrTokenGeneration, "Failed to generate token", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("ConnectWallet: failed to generate token", zap.Error(err))
		return
	}

	logger.Logger.Info("ConnectWallet: user connected", zap.Uint("user_id", user.ID), zap.String("wallet_address", user.WalletAddress))
	c.JSON(http.StatusOK, ConnectWalletResponse{
		Message: "User connected",
		User:    user,
		Token:   token,
	})
}

// GetProfileResponse 用户资料的响应体
type GetProfileResponse struct {
	// User 用户信息
	User models.User `json:"user"`
}

// GetProfile godoc
// @Summary 获取用户资料
// @Description 获取当前认证用户的详细资料。
// @Tags 用户
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} GetProfileResponse "用户资料"
// @Failure 401 {object} errors.APIError "未授权"
// @Failure 500 {object} errors.APIError "服务器错误"
// @Router /api/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	// 从上下文中获取用户ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		apiErr := errors.NewAPIError(errors.ErrInternal, "Failed to get user ID")
		c.Error(apiErr)
		logger.Logger.Error("GetProfile: userID not found in context")
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		apiErr := errors.NewAPIError(errors.ErrInternal, "Invalid user ID format")
		c.Error(apiErr)
		logger.Logger.Error("GetProfile: userID format incorrect", zap.Any("userID", userIDInterface))
		return
	}

	var user models.User
	result := h.DB.First(&user, userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			apiErr := errors.NewAPIError(errors.ErrNotFound, "User not found")
			c.Error(apiErr)
			logger.Logger.Warn("GetProfile: user not found", zap.Uint("user_id", userID))
			return
		}
		apiErr := errors.NewAPIError(errors.ErrDatabase, "Database query error", result.Error.Error())
		c.Error(apiErr)
		logger.Logger.Error("GetProfile: database query error", zap.Error(result.Error))
		return
	}

	logger.Logger.Info("GetProfile: profile retrieved", zap.Uint("user_id", user.ID))
	c.JSON(http.StatusOK, GetProfileResponse{
		User: user,
	})
}
