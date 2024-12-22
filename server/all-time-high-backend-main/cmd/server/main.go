// cmd/server/main.go

// @title Go Web Backend API
// @version 1.0
// @description 这是一个使用Go、Gin、Gorm和PostgreSQL构建的Web后端项目，支持通过Phantom钱包连接用户，并实现了JWT认证。
// @termsOfService http://your-terms-of-service-url.com

// @contact.name API Support
// @contact.url http://www.yourwebsite.com/support
// @contact.email support@yourwebsite.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host api-test.all-time-high.ai

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"log"

	"go.uber.org/zap"

	_ "github.com/GabbyWorld/all-time-high-backend/docs" // 导入生成的Swagger文档
	"github.com/GabbyWorld/all-time-high-backend/internal/config"
	"github.com/GabbyWorld/all-time-high-backend/internal/logger"
	"github.com/GabbyWorld/all-time-high-backend/internal/repository"
	"github.com/GabbyWorld/all-time-high-backend/internal/router"
	"github.com/GabbyWorld/all-time-high-backend/pkg/utils"
	"github.com/joho/godotenv"
)

// main 是应用程序的入口点
func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading configuration from environment variables")
	}
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化日志
	if err := logger.InitLogger(
		cfg.Logger.Level,
		cfg.Logger.LogFile,
		cfg.Logger.MaxSize,
		cfg.Logger.MaxBackups,
		cfg.Logger.MaxAge,
		cfg.Logger.Compress,
	); err != nil {
		log.Fatalf("Could not initialize logger: %v", err)
	}
	defer logger.SyncLogger()

	// 连接数据库并自动迁移
	repo, err := repository.NewRepository(cfg)
	if err != nil {
		logger.Logger.Fatal("Could not connect to the database", zap.Error(err))
	}

	// 初始化JWTManager
	jwtManager, err := utils.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Expiration)
	if err != nil {
		logger.Logger.Fatal("Could not initialize JWT Manager", zap.Error(err))
	}

	// 打印数据库连接状态（可选）
	logger.Logger.Info("Database connection established and migrations run")

	// 设置路由，并传递数据库实例和JWTManager
	r := router.SetupRouter(repo.DB, jwtManager, cfg)

	// 启动服务器
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		logger.Logger.Fatal("Could not run the server", zap.Error(err))
	}
}
