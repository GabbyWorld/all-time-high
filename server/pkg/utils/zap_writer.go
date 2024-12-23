package utils

import (
	"go.uber.org/zap"
)

// ZapWriter 适配 gorm/logger.Writer 接口
type ZapWriter struct {
	Logger *zap.Logger
}

// Printf 实现 gorm/logger.Writer 接口
func (z *ZapWriter) Printf(format string, args ...interface{}) {
	z.Logger.Sugar().Infof(format, args...)
}
