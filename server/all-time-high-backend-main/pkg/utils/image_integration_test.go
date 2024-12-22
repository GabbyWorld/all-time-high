package utils

import (
	"testing"

	"github.com/GabbyWorld/all-time-high-backend/internal/config"
	"github.com/GabbyWorld/all-time-high-backend/internal/logger"
	"github.com/joho/godotenv"
)

func TestGenerateImage_RealAPI(t *testing.T) {
	// 加载 .env 文件
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
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
		t.Fatalf("Could not initialize logger: %v", err)
	}
	defer logger.SyncLogger()

	prompt := " Pixel art of a futuristic sci-fi hero character for a 2D game, Elon Musk, holding a robotic Shiba Inu dog, 32x32 pixel size, dark green matrix theme, clear top-down perspective, no angled view, isolated on transparent background, suitable for vertical scrolling shooter"

	additionalParams := map[string]string{
		"model":         cfg.StableDiffusion.DefaultModel,
		"output_format": cfg.StableDiffusion.OutputFormat,
		"aspect_ratio":  cfg.StableDiffusion.AspectRatio,
		"cfg_scale":     cfg.StableDiffusion.Scale, // 可根据需求调整
		"seed":          "0",                       // 0表示随机种子
	}
	imageBytes, err := GenerateImage(cfg, prompt, additionalParams)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(imageBytes) == 0 {
		t.Error("Expected non-empty image URL")
	}

	imageS3URL, err := UploadImageToS3(cfg, imageBytes)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	t.Logf("Generated Image URL: %s", imageS3URL)
}
