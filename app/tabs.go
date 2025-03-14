package app

import (
	"database/sql"
)

type Tabs struct {
	DB *sql.DB
}

func NewTabs(db *sql.DB) *Tabs {
	return &Tabs{
		DB: db,
	}
}
