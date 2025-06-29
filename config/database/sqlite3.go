package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // SQLite3 driver
)

type Sqlite3 struct {
	DB *sql.DB
}

func getDBPath() string {
	exe, err := os.Executable()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	exeDir := filepath.Dir(exe)
	return filepath.Join(exeDir, "app.db")
}

func NewSqlite3DB() *Sqlite3 {
	// Connect to the SQLite database
	db, err := sql.Open("sqlite3", getDBPath())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the SQLite database successfully.")

	return &Sqlite3{
		DB: db,
	}
}

func (s *Sqlite3) Conn() *sql.DB {
	return s.DB
}

func (s *Sqlite3) CloseConn() {
	s.DB.Close()
}
