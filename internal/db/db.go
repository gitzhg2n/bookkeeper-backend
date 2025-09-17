package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"path/filepath"

	"bookkeeper-backend/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Initialize(cfg *config.Config) (*sql.DB, *gorm.DB, error) {
	dbPath := cfg.DatabaseURL
	if dbPath == "" {
		dbPath = "bookkeeper.db"
	}
	_ = ensureDir(filepath.Dir(dbPath))

	// First, create the SQL connection for migrations
	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("open sqlite connection: %w", err)
	}

	// Run migrations using the new migration system
	logger := slog.Default()
	if err := RunMigrations(sqlDB, logger); err != nil {
		return nil, nil, fmt.Errorf("migrations: %w", err)
	}

	// Now create GORM connection
	gormDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		sqlDB.Close()
		return nil, nil, fmt.Errorf("open gorm sqlite: %w", err)
	}

	// Verify GORM can access the database
	gormSqlDB, err := gormDB.DB()
	if err != nil {
		sqlDB.Close()
		return nil, nil, fmt.Errorf("gorm sql db: %w", err)
	}

	// Close the original connection and return GORM's connection
	sqlDB.Close()
	return gormSqlDB, gormDB, nil
}

func ensureDir(_ string) error { return nil }
