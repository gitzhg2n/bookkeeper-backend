package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Port                string
	DatabaseURL         string
	JWTSecret           []byte
	LogLevel            string
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	IdleTimeout         time.Duration
	ShutdownTimeout     time.Duration
	DatabaseTimeout     time.Duration
	RateLimitWindow     time.Duration
	RateLimitMaxRequests int
}

// Load loads configuration from environment variables with sensible defaults
func Load() *Config {
	config := &Config{
		Port:                getEnv("PORT", "3000"),
		DatabaseURL:         getEnv("DATABASE_URL", "bookkeeper.db"),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
		ReadTimeout:         parseDuration("READ_TIMEOUT", "15s"),
		WriteTimeout:        parseDuration("WRITE_TIMEOUT", "15s"),
		IdleTimeout:         parseDuration("IDLE_TIMEOUT", "60s"),
		ShutdownTimeout:     parseDuration("SHUTDOWN_TIMEOUT", "30s"),
		DatabaseTimeout:     parseDuration("DATABASE_TIMEOUT", "5s"),
		RateLimitWindow:     parseDuration("RATE_LIMIT_WINDOW", "1m"),
		RateLimitMaxRequests: parseInt("RATE_LIMIT_MAX_REQUESTS", 100),
	}

	// JWT Secret validation
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters long")
	}
	config.JWTSecret = []byte(jwtSecret)

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(key, defaultValue string) time.Duration {
	value := getEnv(key, defaultValue)
	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("Invalid duration for %s: %s, using default %s", key, value, defaultValue)
		duration, _ = time.ParseDuration(defaultValue)
	}
	return duration
}

func parseInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Invalid integer for %s: %s, using default %d", key, value, defaultValue)
		return defaultValue
	}
	return intValue
}