package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func userDataDir(appName string) (string, error) {
	base, err := os.UserConfigDir()
	if err != nil || base == "" {
		return "", fmt.Errorf("cannot resolve user config dir: %w", err)
	}
	return filepath.Join(base, appName), nil
}

func ensureDB(appName string) (string, error) {
	dir, err := userDataDir(appName)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	// If you omitted runtimeDBName in the embed files, default here:
	if runtimeDBName == "" {
		runtimeDBName = "app.db"
	}
	dst := filepath.Join(dir, runtimeDBName)

	if _, statErr := os.Stat(dst); errors.Is(statErr, os.ErrNotExist) {
		if len(embeddedDB) == 0 {
			return "", fmt.Errorf("embedded database is empty (check go:embed path and build tags)")
		}
		if writeErr := os.WriteFile(dst, embeddedDB, 0o644); writeErr != nil {
			return "", writeErr
		}
	}
	return dst, nil
}

type Sqlite3 struct{ DB *sql.DB }

func NewSqlite3DB(ctx context.Context) (*Sqlite3, error) {
	dbPath, err := ensureDB("dbmx")
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to the SQLite database successfully.")
	return &Sqlite3{DB: db}, nil
}

func (s *Sqlite3) Conn() *sql.DB { return s.DB }
func (s *Sqlite3) CloseConn()    { _ = s.DB.Close() }
