package db

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"os"
	"path/filepath"
	"runtime"
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

//go:embed schema.sql
var schema string

func Migrate(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, schema)
	if err != nil {
		return err
	}
	return nil
}
