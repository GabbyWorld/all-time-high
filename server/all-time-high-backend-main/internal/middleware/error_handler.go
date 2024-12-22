// internal/middleware/error_handler.go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/GabbyWorld/all-time-high-backend/internal/errors"
	"github.com/GabbyWorld/all-time-high-backend/internal/logger"
)

// ErrorHandlerMiddleware 捕获未处理的错误并返回统一的错误响应
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用 defer 捕获panic
		defer func() {
			if r := recover(); r != nil {
				// 记录panic
				logger.Logger.Error("Panic recovered", zap.Any("error", r))

				// 返回内部服务器错误
				apiErr := errors.NewAPIError(errors.ErrInternal, "Internal server error")
				errors.RespondWithError(c, http.StatusInternalServerError, apiErr)
			}
		}()

		// 继续处理请求
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			// 获取第一个错误
			err := c.Errors[0].Err

			// 根据错误类型确定响应
			switch e := err.(type) {
			case *errors.APIError:
				// 根据ErrorCode决定HTTP状态码
				var statusCode int
				switch e.Code {
				case errors.ErrInvalidRequest, errors.ErrValidation:
					statusCode = http.StatusBadRequest
				case errors.ErrUnauthorized:
					statusCode = http.StatusUnauthorized
				case errors.ErrForbidden:
					statusCode = http.StatusForbidden
				case errors.ErrNotFound:
					statusCode = http.StatusNotFound
				default:
					statusCode = http.StatusInternalServerError
				}
				errors.RespondWithError(c, statusCode, e)
				logger.Logger.Warn("ErrorHandlerMiddleware: API error", zap.String("code", string(e.Code)), zap.String("message", e.Message))
			default:
				// 未知错误，返回内部服务器错误
				apiErr := errors.NewAPIError(errors.ErrInternal, "Internal server error")
				errors.RespondWithError(c, http.StatusInternalServerError, apiErr)
				logger.Logger.Error("ErrorHandlerMiddleware: unknown error", zap.Any("error", err))
			}
		}
	}
}
