package model

// Connection represents a row in the connections table in sqlite3 which represents a connection to a database engine server
type Connection struct {
	ID          int64
	Engine      string
	Host        string
	Port        string
	Username    string
	Password    string
	Database    string
	Name        string
	Env         string
	Color       string
	IsAdvanced  bool
	SSLMode     string
	ClientKey   []byte
	ClientCert  []byte
	RootCACert  []byte
	OverSSH     bool
	SSHHost     string
	SSHPort     string
	SSHUsername string
	SSHPassword string
	UseSSHKey   bool
	SSHKey      []byte

	// Only set for active connection

	IsActive bool
}

type ConnectionTable struct {
	ID       int64
	Name     string
	Env      string
	Engine   string
	Host     string
	Database string
}

type Database struct {
	// ID uniquely identifies the active database within a connection
	ID string
	// PostgresConnectionID is the primary key id of the postgres connection in the sqlite3 database
	ConnectionID   int64
	ConnectionName string
	Name           string
	Color          string

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
	TotalRows    int64    `json:"totalRows"`
	RowsAffected int64    `json:"rowsAffected"`
	Message      string   `json:"message"`

	// If query output contains data of only one table and output also contains id primary key, its name will be stored here
	// Else it will be empty
	TableName string `json:"tableName"`

	// ExecutionTime is the time taken to execute the query in milliseconds
	ExecutionTime int64 `json:"executionTime"`
}

type Output struct {
	Columns []string `json:"columns"`
	Rows    [][]Cell `json:"rows"`
}

type Structure struct {
	Columns []string `json:"columns"`
	Rows    [][]Cell `json:"rows"`
}

type Indexes struct {
	Columns []string `json:"columns"`
	Rows    [][]Cell `json:"rows"`
}

type Rules struct {
	Columns []string `json:"columns"`
	Rows    [][]Cell `json:"rows"`
}

type TableInfo struct {
	Structure Structure `json:"structure"`
	Indexes   Indexes   `json:"indexes"`
	Rules     Rules     `json:"rules"`
}

type UpdateCell struct {
	CellID     string
	TableName  string
	RowID      int64
	ColumnName string
	Value      any
}

// ColumnMeta describes one column of a table. It carries everything the insert-row
// form needs to render an input for *any* postgres type: Category is the column
// type's pg_type.typcategory, so widget choice is driven by the catalog instead of
// a hardcoded list of type names, and CastType is the type every value is cast to.
type ColumnMeta struct {
	Name string `json:"name"`
	// DataType is format_type with the type modifier, e.g. "character varying(255)" — display only.
	DataType string `json:"dataType"`
	// CastType is format_type without the type modifier, e.g. "character varying".
	// It picks the input widget on the frontend; it is deliberately NOT used to cast
	// inserted values -- see InsertRow for why an explicit cast corrupts bit/character.
	CastType string `json:"castType"`
	// Category is pg_type.typcategory: B bool, N numeric, S string, D datetime,
	// E enum, A array, U user-defined (json, uuid, inet, bytea, ...), etc.
	Category     string `json:"category"`
	IsNullable   bool   `json:"isNullable"`
	HasDefault   bool   `json:"hasDefault"`
	DefaultValue string `json:"defaultValue"`
	// IsReadOnly marks identity-always and generated-stored columns, which cannot be inserted into.
	IsReadOnly   bool     `json:"isReadOnly"`
	IsPrimaryKey bool     `json:"isPrimaryKey"`
	EnumValues   []string `json:"enumValues"`
	Comment      string   `json:"comment"`
}

// InsertValue is one column of a new row. Value nil means SQL NULL; columns the
// user leaves at their database default are omitted from the payload entirely.
type InsertValue struct {
	ColumnName string  `json:"columnName"`
	Value      *string `json:"value"`
}
