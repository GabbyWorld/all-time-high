package utils

import (
	"bytes"
	"fmt"

	"github.com/GabbyWorld/all-time-high-backend/internal/config"
	"github.com/GabbyWorld/all-time-high-backend/internal/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UploadImageToS3 上传图像字节到S3并返回S3 URL
func UploadImageToS3(cfg *config.Config, imageBytes []byte) (string, error) {
	// 创建AWS会话
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.AWS.S3Region),
		Credentials: credentials.NewStaticCredentials(
			cfg.AWS.AccessKeyID,
			cfg.AWS.SecretAccessKey,
			"",
		),
	})
	if err != nil {
		logger.Logger.Error("UploadImageToS3: failed to create AWS session", zap.Error(err))
		return "", fmt.Errorf("failed to create AWS session: %w", err)
	}

	// 创建S3上传器
	uploader := s3manager.NewUploader(sess)

	// 生成唯一的文件名
	fileName := fmt.Sprintf("agents/%s.png", uuid.New().String())

	// 上传到S3
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(cfg.AWS.S3Bucket),
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(imageBytes),
		ContentType: aws.String("image/png"),
	})
	if err != nil {
		logger.Logger.Error("UploadImageToS3: failed to upload to S3", zap.Error(err))
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	// 返回文件的URL
	imageS3URL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", cfg.AWS.S3Bucket, cfg.AWS.S3Region, fileName)
	logger.Logger.Info("UploadImageToS3: image uploaded to S3", zap.String("url", imageS3URL))
	return imageS3URL, nil
}
