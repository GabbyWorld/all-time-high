// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://your-terms-of-service-url.com",
        "contact": {
            "name": "API Support",
            "url": "http://www.yourwebsite.com/support",
            "email": "support@yourwebsite.com"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/agent": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "玩家输入name, ticker, prompt，后端生成description和图片，并保存数据。",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Agent"
                ],
                "summary": "创建Agent",
                "parameters": [
                    {
                        "description": "Agent请求体",
                        "name": "agent",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.AgentRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "创建成功",
                        "schema": {
                            "$ref": "#/definitions/handlers.AgentResponse"
                        }
                    },
                    "400": {
                        "description": "请求参数错误",
                        "schema": {
                            "$ref": "#/definitions/errors.APIError"
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "$ref": "#/definitions/errors.APIError"
                        }
                    },
                    "500": {
                        "description": "服务器错误",
                        "schema": {
                            "$ref": "#/definitions/errors.APIError"
                        }
                    }
                }
            }
        },
        "/api/agent/{id}": {
            "get": {
                "description": "根据 Agent ID 获取对应的 Agent 信息",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Agent"
                ],
                "summary": "获取指定 ID 的 Agent",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Agent ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功返回 Agent",
                        "schema": {
                            "$ref": "#/definitions/handlers.AgentResponse"
                        }
                    },
                    "400": {
                        "description": "请求参数错误",
                        "schema": {
                            "$ref": "#/definitions/errors.APIError"
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "$ref": "#/definitions/errors.APIError"
                        }
                    },
                    "404": {
                        "description": "未找到",
                        "schema": {
                            "$ref": "#/definitions/errors.APIError"
                        }
                    },
                    "500": {
                        "description": "服务器错误",
                        "schema": {
                            "$ref": "#/definitions/errors.APIError"
                        }
                    }
                }
            }
        },
        "/api/agents": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "获取当前登录用户关联的所有Agent记录，并支持分页。",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Agent"
                ],
                "summary": "获取登录用户的所有Agents（分页）",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "页码(默认为1)",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "每页大小(默认为4)",
                        "name": "page_size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功返回所有Agent(包含分页信息)",
                        "schema": {
                            "$ref": "#/definitions/handlers.AgentsResponse"
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "$ref": "#/definitions/errors.APIError"
                        }
                    },
                    "500": {
                        "description": "服务器错误",
                        "schema": {
                            "$ref": "#/definitions/errors.APIError"
                        }
                    }
                }
            }
        },
        "/api/agents/all": {
            "get": {
                "description": "获取数据库中所有Agent记录，并支持分页。",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Agent"
                ],
                "summary": "获取所有Agent（分页）",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "页码(默认为1)",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "每页大小(默认为4)",
                        "name": "page_size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功返回所有Agent(包含分页信息)",
                        "schema": {
                            "$ref": "#/definitions/handlers.AgentsResponse"
                        }
                    },
                    "500": {
                        "description": "服务器错误",
                        "schema": {
                            "$ref": "#/definitions/errors.APIError"
                        }
                    }
                }
            }
        },
        "/api/connect_wallet": {
            "post": {
                "description": "用户通过Phantom钱包连接到系统，创建或更新用户信息，并返回JWT令牌。",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "连接Phantom钱包",
                "parameters": [
                    {
                        "description": "用户连接Phantom钱包请求体",
                        "name": "connect_wallet",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.ConnectWalletRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "用户已连接",
                        "schema": {
                            "$ref": "#/definitions/handlers.ConnectWalletResponse"
                        }
                    },
                    "201": {
                        "description": "用户已创建",
                        "schema": {
                            "$ref": "#/definitions/handlers.ConnectWalletResponse"
                        }
                    },
                    "400": {
                        "description": "请求参数错误",
                        "schema": {
                            "$ref": "#/definitions/errors.APIError"
                        }
                    },
                    "500": {
                        "description": "服务器错误",
                        "schema": {
                            "$ref": "#/definitions/errors.APIError"
                        }
                    }
                }
            }
        },
        "/api/health": {
            "get": {
                "description": "检查服务是否运行正常。",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "健康检查"
                ],
                "summary": "健康检查",
                "responses": {
                    "200": {
                        "description": "服务正常",
                        "schema": {
                            "$ref": "#/definitions/handlers.HealthCheckResponse"
                        }
                    }
                }
            }
        },
        "/api/leaderboard": {
            "get": {
                "description": "获取按照胜利次数、胜率和创建时间排序的前100名 Agent",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Agent"
                ],
                "summary": "获取排行榜",
                "responses": {
                    "200": {
                        "description": "成功返回排行榜",
                        "schema": {
                            "$ref": "#/definitions/handlers.LeaderboardResponse"
                        }
                    },
                    "500": {
                        "description": "服务器错误",
                        "schema": {
                            "$ref": "#/definitions/errors.APIError"
                        }
                    }
                }
            }
        },
        "/api/nonce": {
            "get": {
                "description": "生成一个随机的Nonce并存储到内存，返回给客户端（示例使用内存存储，实际可使用Redis等）",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "生成随机Nonce",
                "responses": {
                    "200": {
                        "description": "返回一个包含nonce字段的JSON对象",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/api/profile": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "获取当前认证用户的详细资料。",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "获取用户资料",
                "responses": {
                    "200": {
                        "description": "用户资料",
                        "schema": {
                            "$ref": "#/definitions/handlers.GetProfileResponse"
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "$ref": "#/definitions/errors.APIError"
                        }
                    },
                    "500": {
                        "description": "服务器错误",
                        "schema": {
                            "$ref": "#/definitions/errors.APIError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "errors.APIError": {
            "type": "object",
            "properties": {
                "code": {
                    "$ref": "#/definitions/errors.ErrorCode"
                },
                "details": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "errors.ErrorCode": {
            "type": "string",
            "enum": [
                "INVALID_REQUEST",
                "UNAUTHORIZED",
                "FORBIDDEN",
                "NOT_FOUND",
                "INTERNAL_ERROR",
                "DATABASE_ERROR",
                "VALIDATION_ERROR",
                "TOKEN_GENERATION_ERROR",
                "TOKEN_VERIFICATION_ERROR"
            ],
            "x-enum-varnames": [
                "ErrInvalidRequest",
                "ErrUnauthorized",
                "ErrForbidden",
                "ErrNotFound",
                "ErrInternal",
                "ErrDatabase",
                "ErrValidation",
                "ErrTokenGeneration",
                "ErrTokenVerification"
            ]
        },
        "handlers.AgentInfo": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "image_url": {
                    "type": "string"
                },
                "market_cap": {
                    "type": "number"
                },
                "name": {
                    "type": "string"
                },
                "ticker": {
                    "type": "string"
                },
                "win_rate": {
                    "type": "number"
                },
                "wins": {
                    "type": "integer"
                }
            }
        },
        "handlers.AgentRequest": {
            "type": "object",
            "required": [
                "name",
                "prompt",
                "ticker"
            ],
            "properties": {
                "name": {
                    "type": "string",
                    "maxLength": 100
                },
                "prompt": {
                    "type": "string"
                },
                "ticker": {
                    "type": "string",
                    "maxLength": 50
                }
            }
        },
        "handlers.AgentResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "image_url": {
                    "type": "string"
                },
                "losses": {
                    "description": "新增",
                    "type": "integer"
                },
                "market_cap": {
                    "type": "number"
                },
                "name": {
                    "type": "string"
                },
                "prompt": {
                    "type": "string"
                },
                "ticker": {
                    "type": "string"
                },
                "token_address": {
                    "type": "string"
                },
                "total": {
                    "description": "新增",
                    "type": "integer"
                },
                "user_wallet_address": {
                    "type": "string"
                },
                "win_rate": {
                    "description": "新增",
                    "type": "number"
                },
                "wins": {
                    "description": "新增",
                    "type": "integer"
                }
            }
        },
        "handlers.AgentsResponse": {
            "type": "object",
            "properties": {
                "agents": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/handlers.AgentResponse"
                    }
                }
            }
        },
        "handlers.ConnectWalletRequest": {
            "type": "object",
            "required": [
                "message",
                "signature",
                "wallet_address"
            ],
            "properties": {
                "message": {
                    "description": "Message 签名消息（必填，用于验证 nonce）",
                    "type": "string"
                },
                "signature": {
                    "description": "Signature 用户签名（必填）",
                    "type": "string"
                },
                "username": {
                    "description": "Username 用户名（可选）",
                    "type": "string",
                    "maxLength": 50
                },
                "wallet_address": {
                    "description": "WalletAddress 即公钥（必填）",
                    "type": "string"
                }
            }
        },
        "handlers.ConnectWalletResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "description": "Message 操作消息",
                    "type": "string"
                },
                "token": {
                    "description": "Token JWT令牌",
                    "type": "string"
                },
                "user": {
                    "description": "User 用户信息",
                    "allOf": [
                        {
                            "$ref": "#/definitions/models.User"
                        }
                    ]
                }
            }
        },
        "handlers.GetProfileResponse": {
            "type": "object",
            "properties": {
                "user": {
                    "description": "User 用户信息",
                    "allOf": [
                        {
                            "$ref": "#/definitions/models.User"
                        }
                    ]
                }
            }
        },
        "handlers.HealthCheckResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "description": "Status 服务状态",
                    "type": "string"
                }
            }
        },
        "handlers.LeaderboardResponse": {
            "type": "object",
            "properties": {
                "leaderboard": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/handlers.AgentInfo"
                    }
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                },
                "wallet_address": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:9100",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Go Web Backend API",
	Description:      "这是一个使用Go、Gin、Gorm和PostgreSQL构建的Web后端项目，支持通过Phantom钱包连接用户，并实现了JWT认证。",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
