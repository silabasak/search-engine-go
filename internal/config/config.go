package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	Providers   ProvidersConfig
	Cache       CacheConfig
	Logging     LoggingConfig
	Security    SecurityConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type ProvidersConfig struct {
	JSONURL      string
	XMLURL       string
	Timeout      time.Duration
	RateLimit    int
}

type CacheConfig struct {
	TTL      time.Duration
	MaxSize  int
}

type LoggingConfig struct {
	Level string
	File  string
}

type SecurityConfig struct {
	JWTSecret           string
	CORSAllowedOrigins  string
}

func Load() (*Config, error) {
	config := &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "search_engine"),
		},
		Providers: ProvidersConfig{
			JSONURL:   getEnv("PROVIDER_JSON_URL", "http://localhost:3001/api/videos"),
			XMLURL:    getEnv("PROVIDER_XML_URL", "http://localhost:3002/api/articles"),
			Timeout:   getEnvAsDuration("PROVIDER_TIMEOUT", 30*time.Second),
			RateLimit: getEnvAsInt("PROVIDER_RATE_LIMIT", 100),
		},
		Cache: CacheConfig{
			TTL:     getEnvAsDuration("CACHE_TTL", 5*time.Minute),
			MaxSize: getEnvAsInt("CACHE_MAX_SIZE", 1000),
		},
		Logging: LoggingConfig{
			Level: getEnv("LOG_LEVEL", "info"),
			File:  getEnv("LOG_FILE", "logs/app.log"),
		},
		Security: SecurityConfig{
			JWTSecret:          getEnv("JWT_SECRET", "default-secret-key"),
			CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:8080"),
		},
	}

	return config, nil
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

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
} 