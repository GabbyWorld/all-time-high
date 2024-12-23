package errors

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ErrorCode 定义了错误代码的类型
type ErrorCode string

// 常见的错误代码
const (
	ErrInvalidRequest    ErrorCode = "INVALID_REQUEST"
	ErrUnauthorized      ErrorCode = "UNAUTHORIZED"
	ErrForbidden         ErrorCode = "FORBIDDEN"
	ErrNotFound          ErrorCode = "NOT_FOUND"
	ErrInternal          ErrorCode = "INTERNAL_ERROR"
	ErrDatabase          ErrorCode = "DATABASE_ERROR"
	ErrValidation        ErrorCode = "VALIDATION_ERROR"
	ErrTokenGeneration   ErrorCode = "TOKEN_GENERATION_ERROR"
	ErrTokenVerification ErrorCode = "TOKEN_VERIFICATION_ERROR"
)

// APIError 定义了API错误的结构
type APIError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
}

// Error 实现了error接口
func (e *APIError) Error() string {
	return e.Message
}

// RespondWithError 统一发送错误响应
func RespondWithError(c *gin.Context, statusCode int, apiErr *APIError) {
	c.JSON(statusCode, gin.H{
		"error": apiErr,
	})
}

// NewAPIError 创建一个新的APIError
func NewAPIError(code ErrorCode, message string, details ...string) *APIError {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}
	return &APIError{
		Code:    code,
		Message: message,
		Details: detail,
	}
}

// GetHTTPStatus 返回对应ErrorCode的HTTP状态码
func GetHTTPStatus(code ErrorCode) int {
	switch code {
	case ErrInvalidRequest, ErrValidation:
		return http.StatusBadRequest
	case ErrUnauthorized, ErrTokenVerification:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrNotFound:
		return http.StatusNotFound
	case ErrInternal, ErrDatabase, ErrTokenGeneration:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
