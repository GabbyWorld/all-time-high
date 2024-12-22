# 使用官方的Go镜像作为构建环境
FROM golang:1.23-alpine AS builder

# 设置环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 创建工作目录
WORKDIR /app

# 复制go.mod和go.sum并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN go build -o all-time-high-backend cmd/server/main.go

# 使用轻量级的Alpine镜像作为运行环境
FROM alpine:latest

# 安装必要的依赖（如ca-certificates）
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 复制构建好的二进制文件
COPY --from=builder /app/all-time-high-backend .

# 赋予执行权限
RUN chmod +x all-time-high-backend

# 创建日志目录
RUN mkdir -p /var/log/all-time-high-backend

# 暴露端口（根据配置文件设置）
EXPOSE 9100

# 设置环境变量（可通过docker run传递）
ENV LOG_FILE=/var/log/all-time-high-backend/app.log

# 启动应用
CMD ["./all-time-high-backend"]