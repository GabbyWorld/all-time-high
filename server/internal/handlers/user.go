// internal/handlers/user.go
package handlers

import (
	"net/http"

	"encoding/base64"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/ed25519"
	"gorm.io/gorm"
	"sync"
	"time"

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
	// WalletAddress 即公钥（必填）
	WalletAddress string `json:"wallet_address" binding:"required"`
	// Username 用户名（可选）
	Username string `json:"username" binding:"omitempty,max=50"`
	// Signature 用户签名（必填）
	Signature string `json:"signature" binding:"required"`
	// Message 签名消息（必填，用于验证 nonce）
	Message string `json:"message" binding:"required"`
}

var (
	nonceStore = make(map[string]time.Time)
	mutex      sync.Mutex
)

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
		apiErr := errors.NewAPIError(errors.ErrValidation, "请求参数验证失败", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("ConnectWallet: validation failed", zap.Error(err))
		return
	}

	// 1. 解码签名
	signatureBytes, err := base64.StdEncoding.DecodeString(req.Signature)
	if err != nil {
		apiErr := errors.NewAPIError(errors.ErrValidation, "无效的签名编码", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("ConnectWallet: invalid signature encoding", zap.Error(err))
		return
	}

	// 2. 解析公钥
	publicKey, err := solana.PublicKeyFromBase58(req.WalletAddress)
	if err != nil {
		apiErr := errors.NewAPIError(errors.ErrValidation, "无效的公钥格式", err.Error())
		c.Error(apiErr)
		logger.Logger.Error("ConnectWallet: invalid public key format", zap.Error(err))
		return
	}

	// 3. 验证签名
	valid := ed25519.Verify(publicKey.Bytes(), []byte(req.Message), signatureBytes)
	if !valid {
		apiErr := errors.NewAPIError(errors.ErrValidation, "签名验证失败")
		c.Error(apiErr)
		logger.Logger.Error("ConnectWallet: signature verification failed")
		return
	}

	// 4. 若使用 Nonce，解析并验证 nonce（示例：从消息中提取）
	var nonce string
	_, err = fmt.Sscanf(req.Message, "Login request: %s", &nonce)
	if err == nil && nonce != "" {
		mutex.Lock()
		timestamp, exists := nonceStore[nonce]
		if exists {
			// 检查 Nonce 是否在 5 分钟有效期内
			if time.Since(timestamp) > 5*time.Minute {
				delete(nonceStore, nonce)
				mutex.Unlock()
				apiErr := errors.NewAPIError(errors.ErrValidation, "Nonce 已过期")
				c.Error(apiErr)
				logger.Logger.Error("ConnectWallet: nonce expired")
				return
			}
			// 验证通过后删除已使用的 Nonce
			delete(nonceStore, nonce)
		}
		mutex.Unlock()

		if !exists {
			apiErr := errors.NewAPIError(errors.ErrValidation, "无效的 Nonce")
			c.Error(apiErr)
			logger.Logger.Error("ConnectWallet: invalid nonce")
			return
		}
	}

	// 签名及 nonce 验证通过后，继续原有逻辑
	var user models.User
	result := h.DB.Where("wallet_address = ?", req.WalletAddress).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// 创建新用户
			user = models.User{
				WalletAddress: req.WalletAddress,
				Username:      req.Username,
			}
			if err := h.DB.Create(&user).Error; err != nil {
				apiErr := errors.NewAPIError(errors.ErrDatabase, "创建用户失败", err.Error())
				c.Error(apiErr)
				logger.Logger.Error("ConnectWallet: failed to create user", zap.Error(err))
				return
			}
			// 生成JWT
			token, err := h.JWTManager.GenerateToken(&user)
			if err != nil {
				apiErr := errors.NewAPIError(errors.ErrTokenGeneration, "生成 token 失败", err.Error())
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
		apiErr := errors.NewAPIError(errors.ErrDatabase, "数据库查询错误", result.Error.Error())
		c.Error(apiErr)
		logger.Logger.Error("ConnectWallet: database query error", zap.Error(result.Error))
		return
	}

	// 如果提供了新的用户名则更新
	if req.Username != "" && req.Username != user.Username {
		user.Username = req.Username
		if err := h.DB.Save(&user).Error; err != nil {
			apiErr := errors.NewAPIError(errors.ErrDatabase, "更新用户失败", err.Error())
			c.Error(apiErr)
			logger.Logger.Error("ConnectWallet: failed to update user", zap.Error(err))
			return
		}
		logger.Logger.Info("ConnectWallet: user updated", zap.Uint("user_id", user.ID), zap.String("username", user.Username))
	}

	// 生成JWT
	token, err := h.JWTManager.GenerateToken(&user)
	if err != nil {
		apiErr := errors.NewAPIError(errors.ErrTokenGeneration, "生成 token 失败", err.Error())
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

// GenerateNonce godoc
// @Summary 生成随机Nonce
// @Description 生成一个随机的Nonce并存储到内存，返回给客户端（示例使用内存存储，实际可使用Redis等）
// @Tags 用户
// @Produce json
// @Success 200 {object} map[string]string "返回一个包含nonce字段的JSON对象"
// @Router /api/generate_nonce [get]
func (h *UserHandler) GenerateNonce(c *gin.Context) {
	nonce := generateRandomNonce()
	mutex.Lock()
	nonceStore[nonce] = time.Now()
	mutex.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"nonce": nonce,
	})
}

func generateRandomNonce() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
