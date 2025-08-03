package model

// PostgresConnection represents a row in the postgres table in sqlite3 which represents a connection to a PostgreSQL server
type PostgresConnection struct {
	ID       int64
	Name     string
	Host     string
	Port     string
	Username string
	Password string
	Database string
	Env      string
	Colour   string
	IsActive bool
}

type Database struct {
	// ID uniquely identifies the active database within a connection
	ID string
	// PostgresConnectionID is the primary key id of the postgres connection in the sqlite3 database
	PostgresConnectionID   int64
	PostgresConnectionName string
	Name                   string
	Colour                 string

	// Only set for active connection

	// Active pool id
	PoolID   string
	IsActive bool

	// Tables and columns are set for the active database
	Tables  []string
	Columns []string
}

type Cell struct {
	Column string `json:"column"`
	Value  string `json:"value"`
}

type QueryResult struct {
	OK           bool     `json:"ok"`
	Columns      []string `json:"columns"`
	Rows         [][]Cell `json:"rows"`
	RowsAffected int64    `json:"rowsAffected"`
	Message      string   `json:"message"`
}

type Output struct {
	Columns []string `json:"columns"`
	Rows    [][]Cell `json:"rows"`
}
