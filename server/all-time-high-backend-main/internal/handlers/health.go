// internal/handlers/health.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheckResponse 健康检查的响应体
type HealthCheckResponse struct {
	// Status 服务状态
	Status string `json:"status"`
}

// HealthCheck godoc
// @Summary 健康检查
// @Description 检查服务是否运行正常。
// @Tags 健康检查
// @Accept  json
// @Produce  json
// @Success 200 {object} HealthCheckResponse "服务正常"
// @Router /api/health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthCheckResponse{
		Status: "UP4",
	})
}
