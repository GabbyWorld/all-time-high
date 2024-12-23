package router

import (
	"time"

	"github.com/GabbyWorld/all-time-high-backend/internal/config"
	"github.com/GabbyWorld/all-time-high-backend/internal/handlers"
	"github.com/GabbyWorld/all-time-high-backend/internal/logger"
	"github.com/GabbyWorld/all-time-high-backend/internal/middleware"
	"github.com/GabbyWorld/all-time-high-backend/pkg/utils"
	"github.com/gin-contrib/cors"
	zap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, jwtManager *utils.JWTManager, cfg *config.Config) *gin.Engine {
	r := gin.New()

	// 添加Zap日志中间件
	r.Use(zap.RecoveryWithZap(logger.Logger, true))

	// 集成统一错误处理中间件
	r.Use(middleware.ErrorHandlerMiddleware())

	// 配置CORS中间件
	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     cfg.CORS.AllowedMethods,
		AllowHeaders:     cfg.CORS.AllowedHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           time.Duration(cfg.CORS.MaxAge) * time.Second,
	}

	// 应用CORS中间件
	r.Use(cors.New(corsConfig))

	// Swagger路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 用户相关路由
	userHandler := handlers.UserHandler{
		DB:         db,
		JWTManager: jwtManager,
	}

	// Agent相关路由
	agentHandler := handlers.AgentHandler{
		DB:         db,
		Config:     cfg,
		JWTManager: jwtManager,
	}

	agentWSHandler := &handlers.AgentWebSocketHandler{
		DB: db,
	}

	battleWSHandler := handlers.NewBattleWebSocketHandler(db)
	battleService := handlers.NewBattleService(db, battleWSHandler, cfg)
	battleService.StartPriceMonitoring()

	api := r.Group("/api")
	{
		// 健康检查路由
		api.GET("/health", handlers.HealthCheck)
		api.POST("/connect_wallet", userHandler.ConnectWallet)

		// WebSocket路由（无需JWT认证，示例可根据需要调整认证逻辑）
		api.GET("/ws/agents", agentWSHandler.HandleAgentWebSocket)
		api.GET("/ws/battle", battleWSHandler.HandleBattleWebSocket)

		// 受保护的路由组
		protected := api.Group("/")
		protected.Use(middleware.JWTAuthMiddleware(jwtManager))
		{
			protected.GET("/profile", userHandler.GetProfile)
			protected.POST("/agent", agentHandler.CreateAgent)   // 新增Agent路由
			protected.GET("/agents", agentHandler.GetUserAgents) // 新增Agent查询路由
			protected.GET("/agents/all", agentHandler.GetAllAgents)
			protected.GET("/battle", battleService.GetBattle)
			protected.GET("/battles", battleService.GetBattles)
		}
	}

	return r
}
