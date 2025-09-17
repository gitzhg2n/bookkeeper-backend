package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"bookkeeper-backend/config"
	"bookkeeper-backend/internal/db"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var (
		up    = flag.Bool("up", false, "Run all pending migrations")
		down  = flag.Int("down", 0, "Roll back N migrations")
		force = flag.Int("force", -1, "Force set migration version (use with caution)")
		status = flag.Bool("status", false, "Show current migration status")
	)
	flag.Parse()

	cfg := config.Load()
	logger := setupLogger(cfg)

	dbPath := cfg.DatabaseURL
	if dbPath == "" {
		dbPath = "bookkeeper.db"
	}

	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	switch {
	case *up:
		if err := db.RunMigrations(sqlDB, logger); err != nil {
			logger.Error("migration failed", "error", err)
			os.Exit(1)
		}
		logger.Info("migrations completed successfully")

	case *down > 0:
		if err := db.MigrateDown(sqlDB, *down, logger); err != nil {
			logger.Error("rollback failed", "error", err)
			os.Exit(1)
		}
		logger.Info("rollback completed successfully")

	case *force >= 0:
		if err := db.ForceMigrationVersion(sqlDB, uint(*force), logger); err != nil {
			logger.Error("force version failed", "error", err)
			os.Exit(1)
		}
		logger.Info("migration version forced successfully")

	case *status:
		version, dirty, err := db.GetMigrationVersion(sqlDB)
		if err != nil {
			logger.Error("failed to get migration version", "error", err)
			os.Exit(1)
		}
		logger.Info("migration status", "version", version, "dirty", dirty)
		if dirty {
			logger.Warn("database is in a dirty state - manual intervention may be required")
		}

	default:
		fmt.Println("Usage: migrate [options]")
		fmt.Println("Options:")
		fmt.Println("  -up           Run all pending migrations")
		fmt.Println("  -down N       Roll back N migrations")
		fmt.Println("  -force N      Force set migration version to N (use with caution)")
		fmt.Println("  -status       Show current migration status")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  migrate -up")
		fmt.Println("  migrate -down 1")
		fmt.Println("  migrate -status")
		fmt.Println("  migrate -force 2")
	}
}

func setupLogger(cfg *config.Config) *slog.Logger {
	level := slog.LevelInfo
	switch cfg.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}
