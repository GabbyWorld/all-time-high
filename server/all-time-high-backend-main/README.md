# Go Web Backend

这是一个使用 Go 语言构建的 Web 后端项目，基于 Gin 框架，使用 PostgreSQL 作为数据库，配置管理使用 Viper，并实现了通过 Phantom 钱包连接用户和 JWT 认证机制。

## 运行项目

1. 复制 `.env.example` 为 `.env` 并填写相应配置。
2. 安装依赖：

   ```bash
   make deps
   ```

3. 生成 Swagger 文档（可选）：

   ```bash
   make swagger
   ```

4. 运行服务器：

   ```bash
   make run
   ```

服务器将运行在 `http://localhost:9100`。

## Swagger API 文档

访问 [http://localhost:9100/swagger/index.html](http://localhost:9100/swagger/index.html) 查看 API 文档。

## 目录结构

- `cmd/server`：应用入口
- `internal/config`：配置管理
- `internal/handlers`：HTTP 请求处理器
- `internal/models`：数据模型
- `internal/repository`：数据库访问层
- `internal/router`：路由设置
- `internal/middleware`：中间件
- `internal/utils`：实用工具，如 JWT 管理
- `pkg`：公共库
- `migrations`：数据库迁移文件
- `docs`：Swagger 文档

## Makefile 任务

- `make build`：构建应用程序
- `make run`：构建并运行应用程序
- `make clean`：清理构建产物
- `make test`：运行所有测试
- `make fmt`：格式化代码
- `make lint`：运行静态代码检查
- `make deps`：管理依赖
- `make env`：生成 `.env` 文件
- `make swagger`：生成 Swagger 文档
- `make help`：显示帮助信息

## 测试 API

使用工具如 `curl` 或 Postman 测试 API 端点，或通过 Swagger UI 直接进行测试。

## 贡献

欢迎贡献！ 请提交 Pull Request 或创建 Issue 以报告问题或提出建议。

## 许可证

MIT License. 详见 [LICENSE](LICENSE) 文件。
