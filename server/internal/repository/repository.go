// internal/repository/repository.go
package repository

import (
	"fmt"
	"time"

	"github.com/GabbyWorld/all-time-high-backend/internal/config"
	"github.com/GabbyWorld/all-time-high-backend/pkg/utils"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(cfg *config.Config) (*Repository, error) {
	var err error
	var db *gorm.DB
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
	)

	// 初始化 ZapLogger
	zapLogger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	// 配置 Gorm 的日志级别为 Warn，减少日志量
	newLogger := logger.New(
		&utils.ZapWriter{Logger: zapLogger},
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	retryCount := 5
	retryInterval := time.Second * 2

	for i := 0; i < retryCount; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
		if err == nil {
			// 尝试ping数据库
			sqlDB, err2 := db.DB()
			if err2 == nil && sqlDB.Ping() == nil {
				break
			} else {
				err = fmt.Errorf("ping database failed: %v", err2)
			}
		}

		// 如果连接失败，等待一段时间再试
		time.Sleep(retryInterval)
		retryInterval = retryInterval * 2 // 指数退避
	}

	if err != nil {
		return nil, fmt.Errorf("could not connect to the database after retries: %w", err)
	}

	// 获取底层的 *sql.DB 以配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池参数（可从cfg中读取）
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.Database.ConnMaxIdleTime) * time.Minute)

	// 尝试进行简单的Ping测试，验证连接可用性
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	// 自动迁移模型
	err = db.AutoMigrate(
	// 在此处列出需要迁移的模型，如：
	// &models.User{},
	// &models.Agent{}, // 添加Agent模型
	// &models.Battle{},
	)
	if err != nil {
		return nil, err
	}

	return &Repository{DB: db}, nil
}
