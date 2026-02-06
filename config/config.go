package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 应用配置
type Config struct {
	// API Keys
	FirecrawlAPIKey string
	SerperAPIKey    string
	OpenAIAPIKey    string
	GoogleAPIKey    string
	GoogleEngineID  string
	BingAPIKey      string

	// 服务器配置
	ServerPort string
	GinMode    string

	// 数据库配置
	DBPath string

	// 存储配置
	StoragePath string
	ReportsPath string

	// 搜索配置
	SearchCacheDays   int
	MaxSearchResults  int

	// AI配置
	LLMModel       string
	LLMTemperature float64
	LLMMaxTokens   int
	LLMBaseURL     string // 支持自定义LLM API地址
}

var AppConfig *Config

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("未找到 .env 文件，使用环境变量")
	}

	config := &Config{
		// API Keys
		FirecrawlAPIKey: getEnv("FIRECRAWL_API_KEY", ""),
		SerperAPIKey:    getEnv("SERPER_API_KEY", ""),
		OpenAIAPIKey:    getEnv("OPENAI_API_KEY", ""),
		GoogleAPIKey:    getEnv("GOOGLE_API_KEY", ""),
		GoogleEngineID:  getEnv("GOOGLE_SEARCH_ENGINE_ID", ""),
		BingAPIKey:      getEnv("BING_API_KEY", ""),

		// 服务器配置
		ServerPort: getEnv("SERVER_PORT", "8080"),
		GinMode:    getEnv("GIN_MODE", "release"),

		// 数据库配置
		DBPath: getEnv("DB_PATH", "./data/competitive.db"),

		// 存储配置
		StoragePath: getEnv("STORAGE_PATH", "./storage"),
		ReportsPath: getEnv("REPORTS_PATH", "./reports"),

		// 搜索配置
		SearchCacheDays:  getEnvAsInt("SEARCH_CACHE_DAYS", 7),
		MaxSearchResults: getEnvAsInt("MAX_SEARCH_RESULTS", 10),

		// AI配置
		LLMModel:       getEnv("LLM_MODEL", "gpt-4"),
		LLMTemperature: getEnvAsFloat("LLM_TEMPERATURE", 0.3),
		LLMMaxTokens:   getEnvAsInt("LLM_MAX_TOKENS", 4000),
		LLMBaseURL:     getEnv("OPENAI_BASE_URL", ""), // 自定义API地址（如DeepSeek、Ollama）
	}

	// 创建必要的目录
	os.MkdirAll(config.StoragePath, 0755)
	os.MkdirAll(config.ReportsPath, 0755)
	os.MkdirAll("./data", 0755)

	AppConfig = config
	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}
