package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host         string
	Port         string
	Username     string
	Password     string
	Database     string
	Charset      string
	InitStatus   int
	InitMode     string
	MaxIdleConns int
	MaxOpenConns int
}

type JWTConfig struct {
	Secret              string
	AccessTokenExpire   int
	RefreshTokenExpire  int
}

type CORSConfig struct {
	AllowOrigins string
}

func Load() (*Config, error) {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("SERVER_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "3306"),
			Username:     getEnv("DB_USERNAME", "root"),
			Password:     getEnv("DB_PASSWORD", ""),
			Database:     getEnv("DB_DATABASE", "adcms"),
			Charset:      getEnv("DB_CHARSET", "utf8mb4"),
			InitStatus:   getEnvAsInt("DB_INIT_STATUS", 1),
			InitMode:     getEnv("DB_INIT_MODE", "safe"),
			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
		},
		JWT: JWTConfig{
			Secret:             getEnv("JWT_SECRET", "your-secret-key"),
			AccessTokenExpire:  getEnvAsInt("JWT_ACCESS_TOKEN_EXPIRE", 7200),
			RefreshTokenExpire: getEnvAsInt("JWT_REFRESH_TOKEN_EXPIRE", 604800),
		},
		CORS: CORSConfig{
			AllowOrigins: getEnv("CORS_ALLOW_ORIGINS", "http://localhost:5666"),
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
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// GetDSN 返回数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.Charset,
	)
}
