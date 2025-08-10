package db

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"bookkeeper-backend/config"
	"bookkeeper-backend/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Initialize(cfg *config.Config) (*sql.DB, *gorm.DB, error) {
	// For early MVP we focus on sqlite file path
	dbPath := cfg.DatabaseURL
	if dbPath == "" {
		dbPath = "bookkeeper.db"
	}

	// Ensure relative path directory exists
	dir := filepath.Dir(dbPath)
	if dir != "." && dir != "" {
		// swallow error if exists
		_ = ensureDir(dir)
	}

	gormDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("open sqlite: %w", err)
	}

	// AutoMigrate only for early development; will switch to pure SQL migrations soon
	if err := gormDB.AutoMigrate(
		&models.User{},
	); err != nil {
		return nil, nil, fmt.Errorf("automigrate: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("sql db: %w", err)
	}
	return sqlDB, gormDB, nil
}

func ensureDir(path string) error {
	return nil // In minimal version we skip; placeholder if we add FS ops
}

func HealthCheck(gdb *gorm.DB) error {
	sqlDB, err := gdb.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func Must(gdb *gorm.DB, err error) *gorm.DB {
	if err != nil {
		log.Fatalf("db error: %v", err)
	}
	return gdb
}