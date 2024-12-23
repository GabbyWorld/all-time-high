# Makefile for Go Web Backend Project

# 项目名称
APP_NAME := all-time-high-backend

# Go 命令
GO := go

# 交叉编译设置
TARGET_OS := linux
TARGET_ARCH := amd64

# 默认目标
.PHONY: all
all: build

# 构建应用程序
.PHONY: build
build:
	@echo "Building the application for $(TARGET_OS)/$(TARGET_ARCH)..."
	mkdir -p bin
	CGO_ENABLED=0 GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) $(GO) build -o bin/$(APP_NAME) ./cmd/server

# 运行应用程序 (仅适用于本地运行)
.PHONY: run
run:
	@echo "Running the application..."
	go run ./cmd/server/main.go

# 清理构建产物
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	$(GO) clean
	rm -rf bin/
	rm -f config.env

# 运行测试
.PHONY: test
test:
	@echo "Running tests..."
	$(GO) test ./... -v

# 格式化代码
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# 运行静态代码检查
.PHONY: lint
lint:
	@echo "Running linter..."
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run

# 获取依赖
.PHONY: deps
deps:
	@echo "Tidying up dependencies..."
	$(GO) mod tidy

# 生成 .env 文件示例
.PHONY: env
env:
	@echo "Creating .env file from .env.example..."
	cp .env.example .env

# 生成Swagger文档
.PHONY: swagger
swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/server/main.go

# 帮助信息
.PHONY: help
help:
	@echo "Makefile commands:"
	@echo "  make            - 默认目标，构建应用程序"
	@echo "  make build      - 构建应用程序 (交叉编译为 $(TARGET_OS)/$(TARGET_ARCH))"
	@echo "  make run        - 构建并运行应用程序"
	@echo "  make clean      - 清理构建产物"
	@echo "  make test       - 运行所有测试"
	@echo "  make fmt        - 格式化代码"
	@echo "  make lint       - 运行静态代码检查"
	@echo "  make deps       - 管理依赖"
	@echo "  make env        - 生成 .env 文件"
	@echo "  make swagger    - 生成Swagger文档"
	@echo "  make help       - 显示帮助信息"
