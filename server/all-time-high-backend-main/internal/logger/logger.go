package logger

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 是全局的Zap logger实例
var Logger *zap.Logger

// InitLogger 初始化Zap logger
func InitLogger(logLevel string, logFile string, maxSize, maxBackups, maxAge int, compress bool) error {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		return err
	}

	// 确保日志目录存在（如果不是/dev/stdout）
	if logFile != "/dev/stdout" && logFile != "/dev/stderr" {
		logDir := getLogDir(logFile)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}
	}

	var writeSyncer zapcore.WriteSyncer
	if logFile == "/dev/stdout" || logFile == "/dev/stderr" {
		writeSyncer = zapcore.Lock(os.Stdout)
	} else {
		// 配置Lumberjack作为日志输出器
		lumberjackLogger := &lumberjack.Logger{
			Filename:   logFile,    // 日志文件路径
			MaxSize:    maxSize,    // 单个日志文件的最大大小（MB）
			MaxBackups: maxBackups, // 保留旧日志文件的最大数量
			MaxAge:     maxAge,     // 旧日志文件的最大保留天数
			Compress:   compress,   // 是否压缩旧日志文件
		}
		writeSyncer = zapcore.AddSync(lumberjackLogger)
	}

	// 设置Zap的Encoder配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 选择JSON或Console编码器
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	// encoder := zapcore.NewConsoleEncoder(encoderConfig) // 可选：使用控制台编码器

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevelAt(level)

	// 创建核心（Core）
	core := zapcore.NewCore(encoder, writeSyncer, atomicLevel)

	// 构建Logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	Logger = logger
	return nil
}

// getLogDir 从日志文件路径中提取目录部分
func getLogDir(logFile string) string {
	lastSlash := -1
	for i := len(logFile) - 1; i >= 0; i-- {
		if logFile[i] == '/' || logFile[i] == '\\' {
			lastSlash = i
			break
		}
	}
	if lastSlash == -1 {
		return "."
	}
	return logFile[:lastSlash]
}

// SyncLogger 同步日志
func SyncLogger() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}
