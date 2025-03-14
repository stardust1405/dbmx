package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite3 driver
)

type Sqlite3 struct {
	DB *sql.DB
}

func NewSqlite3DB() *Sqlite3 {
	// Connect to the SQLite database
	db, err := sql.Open("sqlite3", "app.db")
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
