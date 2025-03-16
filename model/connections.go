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
	Colour   *string
	IsActive bool
}

type Database struct {
	PostgresConnectionID int64
	PoolID               int64
	Name                 string
	IsActive             bool
	Tables               []string
}

type GenericResponse struct {
	OK           bool          `json:"ok"`
	Data         []interface{} `json:"data"`
	RowsAffected int64         `json:"rowsAffected"`
	Message      string        `json:"message"`
}
