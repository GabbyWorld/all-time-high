// internal/middleware/jwt.go
package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/GabbyWorld/all-time-high-backend/internal/errors"
	"github.com/GabbyWorld/all-time-high-backend/internal/logger"
	"github.com/GabbyWorld/all-time-high-backend/pkg/utils"
)

// JWTAuthMiddleware 创建一个JWT认证中间件
func JWTAuthMiddleware(jwtManager *utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取Authorization字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			apiErr := errors.NewAPIError(errors.ErrUnauthorized, "Missing Authorization header")
			c.Error(apiErr)
			logger.Logger.Warn("JWTAuthMiddleware: missing Authorization header")
			c.Abort()
			return
		}

		// 检查Authorization字段的格式是否为Bearer {token}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			apiErr := errors.NewAPIError(errors.ErrUnauthorized, "Invalid Authorization format")
			c.Error(apiErr)
			logger.Logger.Warn("JWTAuthMiddleware: invalid Authorization format", zap.String("auth_header", authHeader))
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 验证Token
		claims, err := jwtManager.VerifyToken(tokenString)
		if err != nil {
			apiErr := errors.NewAPIError(errors.ErrTokenVerification, "Invalid token", err.Error())
			c.Error(apiErr)
			logger.Logger.Warn("JWTAuthMiddleware: token verification failed", zap.Error(err))
			c.Abort()
			return
		}

		// 将用户ID存储在上下文中，供后续处理器使用
		c.Set("userID", claims.UserID)

		c.Next()
	}
}
