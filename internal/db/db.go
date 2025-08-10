package db

import (
	"database/sql"
	"fmt"
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

	gormDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("open sqlite: %w", err)
	}

	if err := RunMigrations(gormDB); err != nil {
		return nil, nil, fmt.Errorf("migrations: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("sql db: %w", err)
	}
	return sqlDB, gormDB, nil
}

func ensureDir(_ string) error { return nil }