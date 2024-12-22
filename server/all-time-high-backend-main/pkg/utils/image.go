package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/GabbyWorld/all-time-high-backend/internal/config"
	"github.com/GabbyWorld/all-time-high-backend/internal/logger"
	"go.uber.org/zap"
)

// GenerateImageResponse 定义 Stable Diffusion API 的 JSON 响应结构
type GenerateImageResponse struct {
	Artifacts []struct {
		Base64 string `json:"base64"`
		Seed   int    `json:"seed"`
	} `json:"artifacts"`
}

// GenerateImage 使用 Stable Diffusion API 生成图像
func GenerateImage(cfg *config.Config, prompt string, additionalParams map[string]string) ([]byte, error) {

	// 默认 AcceptHeader 为 image/*
	if cfg.StableDiffusion.AcceptHeader == "" {
		cfg.StableDiffusion.AcceptHeader = "image/*"
	}

	// 创建一个缓冲区和 multipart 写入器
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// 添加必需的表单字段
	if err := writer.WriteField("prompt", prompt); err != nil {
		logger.Logger.Error("GenerateImage: failed to write prompt field", zap.Error(err))
		return nil, fmt.Errorf("failed to write prompt field: %w", err)
	}

	if err := writer.WriteField("mode", "text-to-image"); err != nil {
		logger.Logger.Error("GenerateImage: failed to write default mode field", zap.Error(err))
		return nil, fmt.Errorf("failed to write default mode field: %w", err)
	}

	// 添加其他可选字段
	for key, value := range additionalParams {
		if err := writer.WriteField(key, value); err != nil {
			logger.Logger.Error("GenerateImage: failed to write additional field", zap.String("field", key), zap.Error(err))
			return nil, fmt.Errorf("failed to write additional field %s: %w", key, err)
		}
	}

	// 关闭 multipart 写入器以写入结束边界
	if err := writer.Close(); err != nil {
		logger.Logger.Error("GenerateImage: failed to close multipart writer", zap.Error(err))
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", cfg.StableDiffusion.Endpoint, &requestBody)
	if err != nil {
		logger.Logger.Error("GenerateImage: failed to create HTTP request", zap.Error(err))
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// 设置必要的头部
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.StableDiffusion.APIKey))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", cfg.StableDiffusion.AcceptHeader)

	// 创建 HTTP 客户端
	client := &http.Client{
		Timeout: 60 * time.Second, // 根据需要调整超时时间
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		logger.Logger.Error("GenerateImage: failed to send HTTP request", zap.Error(err))
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Logger.Error("GenerateImage: failed to read response body", zap.Error(err))
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		logger.Logger.Error("GenerateImage: received non-OK status", zap.Int("status_code", resp.StatusCode), zap.String("body", string(bodyBytes)))
		return nil, fmt.Errorf("generate image failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// 根据 AcceptHeader 处理响应
	if cfg.StableDiffusion.AcceptHeader == "application/json" {
		var imageResp GenerateImageResponse
		if err := json.Unmarshal(bodyBytes, &imageResp); err != nil {
			logger.Logger.Error("GenerateImage: failed to unmarshal JSON response", zap.Error(err))
			return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
		}

		if len(imageResp.Artifacts) == 0 {
			logger.Logger.Error("GenerateImage: no artifacts found in response")
			return nil, fmt.Errorf("no artifacts found in response")
		}

		// 解码第一个 Base64 编码的图像
		imageData, err := base64.StdEncoding.DecodeString(imageResp.Artifacts[0].Base64)
		if err != nil {
			logger.Logger.Error("GenerateImage: failed to decode base64 image", zap.Error(err))
			return nil, fmt.Errorf("failed to decode base64 image: %w", err)
		}

		return imageData, nil
	} else {
		// 处理 image/* 响应
		return bodyBytes, nil
	}
}
