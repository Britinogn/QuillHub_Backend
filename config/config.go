package config

import (
    "fmt"
    "os"
    "strconv"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    JWT      JWTConfig
    Email    EmailConfig
    Cloudinary CloudinaryConfig
}

type ServerConfig struct {
    Port        string
    Environment string
    FrontendURL string
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

type RedisConfig struct {
    Host string
    Port string
    URL  string
}

type JWTConfig struct {
    Secret    string
    ExpiresIn string
}

type EmailConfig struct {
    User string
    Pass string
}

type CloudinaryConfig struct {
    CloudName string
    APIKey    string
    APISecret string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
    cfg := &Config{
        Server: ServerConfig{
            Port:        getEnv("PORT", "8080"),
            Environment: getEnv("ENVIRONMENT", "development"),
            FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
        },
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnv("DB_PORT", "5432"),
            User:     getEnv("DB_USER", "admin"),
            Password: getEnv("DB_PASSWORD", "password123"),
            DBName:   getEnv("DB_NAME", "databaseName"),
            SSLMode:  getEnv("DB_SSLMODE", "disable"),
        },
        Redis: RedisConfig{
            Host: getEnv("REDIS_HOST", "localhost"),
            Port: getEnv("REDIS_PORT", "6379"),
            URL:  getEnv("REDIS_URL", ""),
        },
        JWT: JWTConfig{
            Secret:    getEnv("JWT_SECRET", ""),
            ExpiresIn: getEnv("JWT_EXPIRES_IN", "24h"),
        },
        Email: EmailConfig{
            User: getEnv("EMAIL_USER", ""),
            Pass: getEnv("EMAIL_PASS", ""),
        },
        Cloudinary: CloudinaryConfig{
            CloudName: getEnv("CLOUDINARY_CLOUD_NAME", ""),
            APIKey:    getEnv("CLOUDINARY_API_KEY", ""),
            APISecret: getEnv("CLOUDINARY_API_SECRET", ""),
        },
    }

    // Validate required fields
    if cfg.JWT.Secret == "" {
        return nil, fmt.Errorf("JWT_SECRET is required")
    }

    return cfg, nil
}

// getEnv reads an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

// getEnvAsInt reads an environment variable as integer or returns default
func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intVal, err := strconv.Atoi(value); err == nil {
            return intVal
        }
    }
    return defaultValue
}