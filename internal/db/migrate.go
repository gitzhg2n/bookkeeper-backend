package db

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type Migration struct {
	Name string
	SQL  string
}

// RunMigrations reads embedded SQL files and executes them in lexical order.
func RunMigrations(gdb *gorm.DB) error {
	files, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	var migs []Migration
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".sql") {
			continue
		}
		b, err := migrationsFS.ReadFile(filepath.Join("migrations", f.Name()))
		if err != nil {
			return fmt.Errorf("read file %s: %w", f.Name(), err)
		}
		migs = append(migs, Migration{Name: f.Name(), SQL: string(b)})
	}
	sort.Slice(migs, func(i, j int) bool { return migs[i].Name < migs[j].Name })

	// simple migrations table
	if err := gdb.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (filename TEXT PRIMARY KEY, applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`).Error; err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	applied := map[string]struct{}{}
	rows, err := gdb.Raw(`SELECT filename FROM schema_migrations`).Rows()
	if err != nil {
		return fmt.Errorf("query schema_migrations: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		_ = rows.Scan(&name)
		applied[name] = struct{}{}
	}

	for _, m := range migs {
		if _, ok := applied[m.Name]; ok {
			continue
		}
		if err := gdb.Exec(m.SQL).Error; err != nil {
			return fmt.Errorf("apply migration %s: %w", m.Name, err)
		}
		if err := gdb.Exec(`INSERT INTO schema_migrations (filename) VALUES (?)`, m.Name).Error; err != nil {
			return fmt.Errorf("record migration %s: %w", m.Name, err)
		}
	}
	return nil
}