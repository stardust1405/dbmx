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

// ColumnType is the postgres type of one column of a result set, resolved from
// the type OID postgres reports in the row description. It is display metadata
// only -- the values themselves still cross the boundary as strings.
type ColumnType struct {
	Name string `json:"name"`
	// DataType is the type as postgres itself prints it, without a modifier:
	// "integer", "character varying", "timestamp with time zone", "text[]".
	DataType string `json:"dataType"`
	// DisplayType is the short name postgres holds the type under -- "int4",
	// "varchar", "timestamptz" -- which is what a grid header has room for.
	// Arrays keep the printed form instead, since "_text" names nothing a
	// reader would recognise where "text[]" does.
	DisplayType string `json:"displayType"`
	// Category is pg_type.typcategory: B bool, N numeric, S string, D datetime,
	// E enum, A array, U user-defined (json, uuid, inet, bytea, ...), etc. The
	// grid picks a column's header icon from it, so a type the app has never
	// heard of still gets the icon its family deserves.
	Category string `json:"category"`

	// The fields below describe the column a value came from, so they are only
	// set for a result column that is a plain reference to a table column. An
	// expression, an aggregate or a literal has no attribute behind it and
	// leaves all of them at their zero value.

	IsNullable   bool `json:"isNullable"`
	IsPrimaryKey bool `json:"isPrimaryKey"`
	// IsCompositeKey marks a primary key that spans more than one column. Such
	// a column identifies a row only together with its siblings, which is worth
	// saying in a grid that lets rows be edited by their key.
	IsCompositeKey bool `json:"isCompositeKey"`
	IsForeignKey   bool `json:"isForeignKey"`
	// ForeignKeyTable is the table a foreign key points at, for the tooltip.
	// A column under more than one foreign key reports the first by name.
	ForeignKeyTable string `json:"foreignKeyTable"`
}

type QueryResult struct {
	OK      bool     `json:"ok"`
	Columns []string `json:"columns"`
	// ColumnTypes is parallel to Columns. It is empty for results that have no
	// table behind them (an error row, a write's "Rows Affected").
	ColumnTypes  []ColumnType `json:"columnTypes"`
	Rows         [][]Cell     `json:"rows"`
	TotalRows    int64        `json:"totalRows"`
	RowsAffected int64        `json:"rowsAffected"`
	Message      string       `json:"message"`

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

// ColumnDefinition is the editable shape of one column, used by both the
// add-column and edit-column forms. It is deliberately separate from ColumnMeta:
// ColumnMeta describes a column so a *value* can be entered into it, while this
// describes a column so the *column itself* can be created or altered.
type ColumnDefinition struct {
	Name string `json:"name"`
	// DataType is the type as written in DDL, modifier included:
	// "character varying(255)", "numeric(10,2)", "text[]", "public.mood".
	DataType   string `json:"dataType"`
	IsNullable bool   `json:"isNullable"`
	// DefaultValue is a raw SQL expression ("now()", "0", "'draft'::text").
	// Empty means the column has no default.
	DefaultValue string `json:"defaultValue"`
	// Identity is "", "BY DEFAULT" or "ALWAYS" for GENERATED ... AS IDENTITY.
	Identity string `json:"identity"`
	// GeneratedExpression is the GENERATED ALWAYS AS (...) STORED expression.
	// Postgres cannot alter one in place, so on edit it is read-only: the form
	// may clear it (which drops the expression) but not rewrite it.
	GeneratedExpression string `json:"generatedExpression"`
	Collation           string `json:"collation"`
	Comment             string `json:"comment"`
	// UsingExpression is the USING clause of ALTER COLUMN ... TYPE, needed when
	// the old and new types have no assignment cast. Ignored when adding.
	UsingExpression string `json:"usingExpression"`
	// IsPrimaryKey is display-only: the constraint itself is managed under Rules.
	IsPrimaryKey bool `json:"isPrimaryKey"`
}

// IndexDefinition is one index, held as the structured parts of a CREATE INDEX
// rather than as statement text, so the form can round-trip an existing index
// through the catalog without parsing pg_get_indexdef output.
type IndexDefinition struct {
	Name     string `json:"name"`
	Method   string `json:"method"`
	IsUnique bool   `json:"isUnique"`
	// Columns are key column expressions exactly as postgres renders them,
	// operator class and ordering included: "email", "lower(name)", "id DESC".
	Columns []string `json:"columns"`
	// Include are non-key payload columns (INCLUDE, postgres 11+).
	Include []string `json:"include"`
	// Where is the partial-index predicate; empty means the index covers every row.
	Where   string `json:"where"`
	Comment string `json:"comment"`
	// IsConstraint marks an index postgres created to back a constraint. Such an
	// index cannot be dropped or recreated on its own; it is edited under Rules.
	IsConstraint bool `json:"isConstraint"`
	IsPrimary    bool `json:"isPrimary"`
}

// ConstraintDefinition is one table constraint. The body is carried as the text
// pg_get_constraintdef renders -- "CHECK ((price > 0))", "FOREIGN KEY (a) REFERENCES
// b(id) ON DELETE CASCADE" -- because that covers every constraint type postgres
// has, including EXCLUDE, without the model growing a branch per type. The form
// composes the body from structured inputs and still lets it be edited by hand.
type ConstraintDefinition struct {
	Name string `json:"name"`
	// Type is PRIMARY KEY, UNIQUE, CHECK, FOREIGN KEY or EXCLUDE.
	Type       string `json:"type"`
	Definition string `json:"definition"`
	Comment    string `json:"comment"`
	// NotValid adds the constraint without checking existing rows (CHECK and
	// FOREIGN KEY only).
	NotValid bool `json:"notValid"`
	// IsValidated reports whether postgres has verified the existing rows.
	IsValidated bool `json:"isValidated"`
}

// TypeOption is one type the column editor offers. Kind distinguishes the
// built-ins from a database's own enums, domains and composites.
type TypeOption struct {
	Name string `json:"name"`
	// Kind is base, enum, domain, range or composite.
	Kind   string `json:"kind"`
	Schema string `json:"schema"`
	// IsCommon marks the handful of types worth surfacing before the full list.
	IsCommon bool `json:"isCommon"`
	// AcceptsModifier reports whether the type takes a modifier -- a length for
	// character varying, a precision and scale for numeric, fractional seconds for
	// the datetime types. It comes from pg_type.typmodin, so a domain or extension
	// type that takes one is recognised alongside the built-ins.
	AcceptsModifier bool `json:"acceptsModifier"`
	// EnumValues are the labels of an enum type, for display next to its name.
	EnumValues []string `json:"enumValues"`
}

// SchemaEditorOptions is everything the three schema forms need to populate their
// selects, fetched in one round trip when a form opens.
type SchemaEditorOptions struct {
	Types        []TypeOption `json:"types"`
	IndexMethods []string     `json:"indexMethods"`
	Collations   []string     `json:"collations"`
	Tables       []string     `json:"tables"`
}
