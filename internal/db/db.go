package db

import (
	"context"
	"database/sql"
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

func getDBPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	var baseDir string
	switch runtime.GOOS {
	case "linux":
		// Matches: ~/.local/share
		baseDir = filepath.Join(home, ".local", "share")

	case "darwin":
		// Matches: ~/Library/Application Support
		baseDir = filepath.Join(home, "Library", "Application Support")

	case "windows":
		// Matches: %LOCALAPPDATA% (typically C:\Users\Username\AppData\Local)
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData != "" {
			baseDir = localAppData
		} else {
			baseDir = filepath.Join(home, "AppData", "Local")
		}
	default:
		// Fallback for other Unix-like systems
		baseDir = filepath.Join(home, ".local", "share")

	}
	return filepath.Join(baseDir, "todoer", "todoer.db"), nil

}

func Init(ctx context.Context) (*sql.DB, error) {
	dbPath, err := getDBPath()

	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		log.Fatal("Failed to create dir", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	timeOutCtx, timeOutcancel := context.WithTimeout(ctx, 5*time.Second)
	defer timeOutcancel()
	if err := db.PingContext(timeOutCtx); err != nil {
		log.Fatal(err)
	}
	return db, nil

}

//go:embed migrations/*.sql
var migrationFiles embed.FS

func Migrate(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY
		)
	`)
	if _, err := db.ExecContext(ctx, "PRAGMA foreign_keys = ON"); err != nil {
		return err
	}
	if err != nil {
		return err
	}

	files, err := fs.Glob(migrationFiles, "migrations/*.sql")

	if err != nil {
		return err
	}

	sort.Strings(files)

	for _, file := range files {
		version := strings.TrimSuffix(strings.TrimPrefix(file, "migrations/"), ".sql")

		var exists bool

		err := db.QueryRowContext(ctx, `
		SELECT EXISTS(
		SELECT 1 FROM schema_migrations WHERE version = ?)
		`, version).Scan(&exists)

		if err != nil {
			return err
		}
		if exists {
			continue
		}

		sqlBytes, err := migrationFiles.ReadFile(file)
		if err != nil {
			return err
		}
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, string(sqlBytes)); err != nil {
			tx.Rollback()
			return fmt.Errorf("%s: %w", version, err)
		}

		if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations(version) VALUES(?)", version); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

	}

	return nil
}
