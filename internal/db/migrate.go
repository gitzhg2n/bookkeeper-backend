package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations applies all pending database migrations using golang-migrate
func RunMigrations(sqlDB *sql.DB, logger *slog.Logger) error {
	driver, err := sqlite3.WithInstance(sqlDB, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %w", err)
	}

	// Get the absolute path to migrations directory
	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		return fmt.Errorf("could not get migrations path: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}
	defer m.Close()

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		logger.Warn("could not get migration version", "error", err)
	} else if err != migrate.ErrNilVersion {
		logger.Info("database migrations completed", "version", version, "dirty", dirty)
	} else {
		logger.Info("no migrations to apply")
	}

	return nil
}

// MigrateDown rolls back the database by specified steps
func MigrateDown(sqlDB *sql.DB, steps int, logger *slog.Logger) error {
	driver, err := sqlite3.WithInstance(sqlDB, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %w", err)
	}

	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		return fmt.Errorf("could not get migrations path: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Steps(-steps); err != nil {
		return fmt.Errorf("could not run down migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		logger.Warn("could not get migration version", "error", err)
	} else if err != migrate.ErrNilVersion {
		logger.Info("database rollback completed", "version", version, "dirty", dirty)
	}

	return nil
}

// GetMigrationVersion returns the current migration version
func GetMigrationVersion(sqlDB *sql.DB) (uint, bool, error) {
	driver, err := sqlite3.WithInstance(sqlDB, &sqlite3.Config{})
	if err != nil {
		return 0, false, fmt.Errorf("could not create migration driver: %w", err)
	}

	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		return 0, false, fmt.Errorf("could not get migrations path: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"sqlite3",
		driver,
	)
	if err != nil {
		return 0, false, fmt.Errorf("could not create migrate instance: %w", err)
	}
	defer m.Close()

	return m.Version()
}

// ForceMigrationVersion sets the migration version without running migrations
// This is useful for fixing dirty migration states
func ForceMigrationVersion(sqlDB *sql.DB, version uint, logger *slog.Logger) error {
	driver, err := sqlite3.WithInstance(sqlDB, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %w", err)
	}

	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		return fmt.Errorf("could not get migrations path: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Force(int(version)); err != nil {
		return fmt.Errorf("could not force migration version: %w", err)
	}

	logger.Info("forced migration version", "version", version)
	return nil
}
