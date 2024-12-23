package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server          ServerConfig
	Database        DatabaseConfig
	JWT             JWTConfig
	Logger          LoggerConfig
	OpenAI          OpenAIConfig
	ImageAPI        ImageAPIConfig
	StableDiffusion StableDiffusionConfig
	AWS             AWSConfig
	CORS            CORSConfig
	Solana          SolanaConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int // 分钟为单位
	ConnMaxIdleTime int // 分钟为单位
}

type JWTConfig struct {
	Secret     string
	Expiration string
}

type LoggerConfig struct {
	Level      string // e.g., "debug", "info", "warn", "error"
	LogFile    string // e.g., "/var/log/all-time-high/app.log" 或 "./logs/app.log"
	MaxSize    int    // 单个日志文件的最大大小（MB）
	MaxBackups int    // 保留旧日志文件的最大数量
	MaxAge     int    // 旧日志文件的最大保留天数
	Compress   bool   // 是否压缩旧日志文件
}

type OpenAIConfig struct {
	APIKey              string
	CompletionsEndpoint string
}

type ImageAPIConfig struct {
	APIKey   string
	Endpoint string
}

type StableDiffusionConfig struct {
	APIKey         string
	Endpoint       string
	DefaultModel   string
	AspectRatio    string
	NegativePrompt string
	OutputFormat   string
	AcceptHeader   string
	Scale          string
}

type AWSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	S3Bucket        string
	S3Region        string
}

type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	AllowedMethods   []string `mapstructure:"CORS_ALLOWED_METHODS"`
	AllowedHeaders   []string `mapstructure:"CORS_ALLOWED_HEADERS"`
	AllowCredentials bool     `mapstructure:"CORS_ALLOW_CREDENTIALS"`
	MaxAge           int      `mapstructure:"CORS_MAX_AGE"`
}

type SolanaConfig struct {
	RPCEndpoint      string `mapstructure:"SOLANA_RPC_ENDPOINT"`
	WSRPCEndpoint    string `mapstructure:"SOLANA_WSRPC_ENDPOINT"`
	SignerPrivateKey string `mapstructure:"SOLANA_SIGNER_PRIVATE_KEY"`
	IPFSURL          string `mapstructure:"SOLANA_IPFS_URL"`
	TradeURL         string `mapstructure:"SOLANA_TRADE_URL"`
	TokenProgramID   string `mapstructure:"SOLANA_TOKEN_PROGRAM_ID"`
}

func LoadConfig() *Config {
	viper.AutomaticEnv()

	// 设置默认值
	viper.SetDefault("SERVER_PORT", "9100")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("JWT_EXPIRE", "24h")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FILE", "./logs/app.log") // 本地默认日志路径
	viper.SetDefault("LOG_MAX_SIZE", 100)          // MB
	viper.SetDefault("LOG_MAX_BACKUPS", 3)
	viper.SetDefault("LOG_MAX_AGE", 28) // 天
	viper.SetDefault("LOG_COMPRESS", true)
	viper.SetDefault("DB_MAX_OPEN_CONNS", 20)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 10)
	viper.SetDefault("DB_CONN_MAX_LIFETIME", 30)
	viper.SetDefault("DB_CONN_MAX_IDLE_TIME", 5)
	viper.SetDefault("OPENAI_API_KEY", "")
	viper.SetDefault("IMAGE_API_KEY", "")
	viper.SetDefault("IMAGE_API_ENDPOINT", "https://api.openai.com/v1/images/generations") // 示例使用OpenAI的DALL-E
	viper.SetDefault("OPENAI_COMPLETIONS_ENDPOINT", "https://api.openai.com/v1/chat/completions")
	viper.SetDefault("AWS_REGION", "us-east-1")
	viper.SetDefault("AWS_ACCESS_KEY_ID", "")
	viper.SetDefault("AWS_SECRET_ACCESS_KEY", "")
	viper.SetDefault("AWS_S3_BUCKET", "")
	viper.SetDefault("AWS_S3_REGION", "us-east-1")
	// CORS 配置默认值
	viper.SetDefault("CORS_ALLOWED_ORIGINS", []string{"*"})
	viper.SetDefault("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("CORS_ALLOWED_HEADERS", []string{"Origin", "Content-Type", "Authorization"})
	viper.SetDefault("CORS_ALLOW_CREDENTIALS", false)
	viper.SetDefault("CORS_MAX_AGE", 86400) // 24小时
	viper.SetDefault("SOLANA_RPC_ENDPOINT", "https://mainnet.helius-rpc.com/?api-key=f77fbc1f-282a-4bd7-99e7-cad253f17a77")
	viper.SetDefault("SOLANA_WSRPC_ENDPOINT", "wss://mainnet.helius-rpc.com/?api-key=f77fbc1f-282a-4bd7-99e7-cad253f17a77")
	viper.SetDefault("SOLANA_SIGNER_PRIVATE_KEY", "")
	viper.SetDefault("SOLANA_IPFS_URL", "https://pump.fun/api/ipfs")
	viper.SetDefault("SOLANA_TRADE_URL", "https://pumpportal.fun/api/trade-local")
	viper.SetDefault("SOLANA_TOKEN_PROGRAM_ID", "6EF8rcorrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No config file found, reading from environment variables")
	}
	config := &Config{
		Server: ServerConfig{
			Port: viper.GetString("SERVER_PORT"),
		},
		Database: DatabaseConfig{
			Host:            viper.GetString("DB_HOST"),
			Port:            viper.GetInt("DB_PORT"),
			User:            viper.GetString("DB_USER"),
			Password:        viper.GetString("DB_PASSWORD"),
			DBName:          viper.GetString("DB_NAME"),
			SSLMode:         viper.GetString("DB_SSLMODE"),
			MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: viper.GetInt("DB_CONN_MAX_LIFETIME"),
			ConnMaxIdleTime: viper.GetInt("DB_CONN_MAX_IDLE_TIME"),
		},
		JWT: JWTConfig{
			Secret:     viper.GetString("JWT_SECRET"),
			Expiration: viper.GetString("JWT_EXPIRE"),
		},
		Logger: LoggerConfig{
			Level:      viper.GetString("LOG_LEVEL"),
			LogFile:    viper.GetString("LOG_FILE"),
			MaxSize:    viper.GetInt("LOG_MAX_SIZE"),
			MaxBackups: viper.GetInt("LOG_MAX_BACKUPS"),
			MaxAge:     viper.GetInt("LOG_MAX_AGE"),
			Compress:   viper.GetBool("LOG_COMPRESS"),
		},
		OpenAI: OpenAIConfig{
			APIKey:              viper.GetString("OPENAI_API_KEY"),
			CompletionsEndpoint: viper.GetString("OPENAI_COMPLETIONS_ENDPOINT"),
		},
		ImageAPI: ImageAPIConfig{
			APIKey:   viper.GetString("IMAGE_API_KEY"),
			Endpoint: viper.GetString("IMAGE_API_ENDPOINT"),
		},
		StableDiffusion: StableDiffusionConfig{
			APIKey:         viper.GetString("STABILITY_API_KEY"),
			Endpoint:       viper.GetString("STABLE_DIFFUSION_ENDPOINT"),
			DefaultModel:   viper.GetString("STABLE_DIFFUSION_DEFAULT_MODEL"),
			AspectRatio:    viper.GetString("STABLE_DIFFUSION_ASPECT_RATIO"),
			NegativePrompt: viper.GetString("STABLE_DIFFUSION_NEGATIVE_PROMPT"),
			OutputFormat:   viper.GetString("STABLE_DIFFUSION_OUTPUT_FORMAT"),
			AcceptHeader:   viper.GetString("STABLE_DIFFUSION_ACCEPT_HEADER"),
			Scale:          viper.GetString("STABLE_DIFFUSION_SCALE"),
		},
		AWS: AWSConfig{
			Region:          viper.GetString("AWS_REGION"),
			AccessKeyID:     viper.GetString("AWS_ACCESS_KEY_ID"),
			SecretAccessKey: viper.GetString("AWS_SECRET_ACCESS_KEY"),
			S3Bucket:        viper.GetString("AWS_S3_BUCKET"),
			S3Region:        viper.GetString("AWS_S3_REGION"),
		},
		CORS: CORSConfig{
			AllowedOrigins:   viper.GetStringSlice("CORS_ALLOWED_ORIGINS"),
			AllowedMethods:   viper.GetStringSlice("CORS_ALLOWED_METHODS"),
			AllowedHeaders:   viper.GetStringSlice("CORS_ALLOWED_HEADERS"),
			AllowCredentials: viper.GetBool("CORS_ALLOW_CREDENTIALS"),
			MaxAge:           viper.GetInt("CORS_MAX_AGE"),
		},
		Solana: SolanaConfig{
			RPCEndpoint:      viper.GetString("SOLANA_RPC_ENDPOINT"),
			WSRPCEndpoint:    viper.GetString("SOLANA_WSRPC_ENDPOINT"),
			SignerPrivateKey: viper.GetString("SOLANA_SIGNER_PRIVATE_KEY"),
			IPFSURL:          viper.GetString("SOLANA_IPFS_URL"),
			TradeURL:         viper.GetString("SOLANA_TRADE_URL"),
			TokenProgramID:   viper.GetString("SOLANA_TOKEN_PROGRAM_ID"),
		},
	}

	// 验证必要的配置项
	if config.Database.Host == "" || config.Database.User == "" || config.Database.Password == "" || config.Database.DBName == "" {
		log.Fatal("Database configuration is incomplete. Please set DB_HOST, DB_USER, DB_PASSWORD, DB_NAME.")
	}
	if config.OpenAI.APIKey == "" {
		log.Fatal("OpenAI API key is required. Please set OPENAI_API_KEY.")
	}
	if config.AWS.AccessKeyID == "" || config.AWS.SecretAccessKey == "" || config.AWS.S3Bucket == "" {
		log.Fatal("AWS credentials and S3 bucket are required. Please set AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and AWS_S3_BUCKET.")
	}
	if config.Solana.SignerPrivateKey == "" {
		log.Fatal("Solana private keys are required. Please set SOLANA_SIGNER_PRIVATE_KEY.")
	}

	log.Printf("Server will run on port: %s", config.Server.Port)
	log.Printf("Connecting to database: %s@%s:%d/%s with SSL mode: %s", config.Database.User, config.Database.Host, config.Database.Port, config.Database.DBName, config.Database.SSLMode)
	log.Printf("CORS Allowed Origins: %v", config.CORS.AllowedOrigins) // 添加日志

	return config
}
