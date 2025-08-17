package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port                 string
	DeploymentMode       string
	DatabaseURL          string
	JWTSecret            []byte
	LogLevel             string
	ReadTimeout          time.Duration
	WriteTimeout         time.Duration
	IdleTimeout          time.Duration
	ShutdownTimeout      time.Duration
	DatabaseTimeout      time.Duration
	RateLimitWindow      time.Duration
	RateLimitMaxRequests int

	AccessTokenTTL        time.Duration
	RefreshTokenTTL       time.Duration
	PasswordMemoryKiB     uint32
	PasswordTime          uint32
	PasswordParallelism   uint8
	PasswordSaltLength    uint32
	PasswordKeyLength     uint32
	EncryptionKeyVersion  int
	AllowInsecurePassword bool
}

func Load() *Config {
	cfg := &Config{
	DeploymentMode:       getEnv("DEPLOYMENT_MODE", "cloud"),
		Port:                 getEnv("PORT", "3000"),
		DatabaseURL:          getEnv("DATABASE_URL", "bookkeeper.db"),
		LogLevel:             getEnv("LOG_LEVEL", "info"),
		ReadTimeout:          parseDuration("READ_TIMEOUT", "15s"),
		WriteTimeout:         parseDuration("WRITE_TIMEOUT", "15s"),
		IdleTimeout:          parseDuration("IDLE_TIMEOUT", "60s"),
		ShutdownTimeout:      parseDuration("SHUTDOWN_TIMEOUT", "30s"),
		DatabaseTimeout:      parseDuration("DATABASE_TIMEOUT", "5s"),
		RateLimitWindow:      parseDuration("RATE_LIMIT_WINDOW", "1m"),
		RateLimitMaxRequests: parseInt("RATE_LIMIT_MAX_REQUESTS", 100),

		AccessTokenTTL:       parseDuration("ACCESS_TOKEN_TTL", "15m"),
		RefreshTokenTTL:      parseDuration("REFRESH_TOKEN_TTL", "720h"),
		PasswordMemoryKiB:    uintEnv("PASSWORD_MEMORY_KIB", 64*1024),
		PasswordTime:         uintEnv("PASSWORD_TIME", 3),
		PasswordParallelism:  uint8Env("PASSWORD_PARALLELISM", 1),
		PasswordSaltLength:   uintEnv("PASSWORD_SALT_LENGTH", 16),
		PasswordKeyLength:    uintEnv("PASSWORD_KEY_LENGTH", 32),
		EncryptionKeyVersion: parseInt("DATA_ENCRYPTION_KDF_PARAMS_VERSION", 1),
		AllowInsecurePassword: boolEnv("ALLOW_INSECURE_PASSWORD", false),
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters")
	}
	cfg.JWTSecret = []byte(jwtSecret)

	return cfg
}

// (helpers unchanged below)
func getEnv(key, def string) string { val := os.Getenv(key); if val != "" { return val }; return def }
func parseDuration(key, def string) time.Duration { v := getEnv(key, def); d, err := time.ParseDuration(v); if err != nil { log.Printf("invalid duration for %s=%s using default %s", key, v, def); d, _ = time.ParseDuration(def) }; return d }
func parseInt(key string, def int) int { v := os.Getenv(key); if v == "" { return def }; i, err := strconv.Atoi(v); if err != nil { log.Printf("invalid int for %s=%s using default %d", key, v, def); return def }; return i }
func uintEnv(key string, def uint32) uint32 { v := os.Getenv(key); if v == "" { return def }; i, err := strconv.ParseUint(v, 10, 32); if err != nil { log.Printf("invalid uint for %s=%s using default %d", key, v, def); return def }; return uint32(i) }
func uint8Env(key string, def uint8) uint8 { v := os.Getenv(key); if v == "" { return def }; i, err := strconv.ParseUint(v, 10, 8); if err != nil { log.Printf("invalid uint8 for %s=%s using default %d", key, v, def); return def }; return uint8(i) }
func boolEnv(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" { return def }
	switch v {
	case "1","true","TRUE","yes","y": return true
	case "0","false","FALSE","no","n": return false
	}
	return def
}