package model

// PostgresConnection represents a connection to a PostgreSQL server
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
	ID                     string
	PostgresConnectionID   int64
	PostgresConnectionName string
	Name                   string
	Colour                 string

	// Only set for active connection
	PoolID   string
	IsActive bool
	Tables   []string
	Columns  []string
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
