package db

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"bookkeeper-backend/config"
	"bookkeeper-backend/internal/models"

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

	// AutoMigrate (temporary)
	if err := gormDB.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.Household{},
		&models.HouseholdMember{},
		&models.Account{},
		&models.Transaction{},
		&models.Category{},
		&models.Budget{},
	); err != nil {
		return nil, nil, fmt.Errorf("automigrate: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("sql db: %w", err)
	}
	return sqlDB, gormDB, nil
}

func ensureDir(_ string) error { return nil }

func HealthCheck(gdb *gorm.DB) error {
	sqlDB, err := gdb.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}