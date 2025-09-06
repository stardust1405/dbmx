package app

import (
	"context"
	"database/sql"
	"dbmx/model"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type Connections struct {
	DB *sql.DB
	PM *PoolManager
}

func NewConnections(db *sql.DB, pm *PoolManager) *Connections {
	return &Connections{
		DB: db,
		PM: pm,
	}
}

func (m *Connections) GetSqlite3Version() string {
	// Get the version of SQLite
	var sqliteVersion string
	err := m.DB.QueryRow("SELECT sqlite_version()").Scan(&sqliteVersion)
	if err != nil {
		log.Fatal(err)
	}

	return sqliteVersion
}

func (m *Connections) GetPostgresConnections() ([]model.PostgresConnection, error) {
	// Get all postgres connections
	var connections []model.PostgresConnection
	rows, err := m.DB.Query("SELECT * FROM postgres")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var connection model.PostgresConnection
		err := rows.Scan(
			&connection.ID,
			&connection.Name,
			&connection.Host,
			&connection.Port,
			&connection.Username,
			&connection.Password,
			&connection.Env,
			&connection.Colour,
			&connection.Database,
		)
		if err != nil {
			return nil, errors.Wrap(err, "unable to read resultant rows into connection variable")
		}
		connections = append(connections, connection)
	}

	if row_err := rows.Err(); row_err != nil {
		return nil, errors.Wrap(row_err, "unable to read rows")
	}

	return connections, nil
}

func (m *Connections) TestConnectPostgres(p model.PostgresConnection) (bool, error) {
	if p.Database == "" {
		p.Database = "postgres"
	}

	// Build connection string using the credentials
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require", p.Username, p.Password, p.Host, p.Port, p.Database)

	// Establish a connection
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return false, err
	}
	defer conn.Close(ctx)

	// Execute a simple query
	var greeting string
	err = conn.QueryRow(ctx, "SELECT 'Connection Successful!' AS success").Scan(&greeting)
	if err != nil {
		return false, errors.Wrap(err, "failed to query database")
	}

	// Print the result
	fmt.Printf("Greeting: %s\n", greeting)

	return true, nil
}

func (c *Connections) AddPostgresConnection(p model.PostgresConnection) (bool, error) {
	if p.Database == "" {
		p.Database = "postgres"
	}

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM postgres WHERE name = ?)`

	// Execute the query
	err := c.DB.QueryRow(query, p.Name).Scan(&exists)
	if err != nil {
		return false, err
	}

	if exists {
		return false, errors.New("Connection name already exists. Please choose a different name")
	}

	insertStatement, err := c.DB.Prepare("INSERT INTO postgres (name, host, port, username, password, env, colour, database) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return false, errors.Wrap(err, "failed to prepare query to insert new connection in postgres")
	}

	_, err = insertStatement.Exec(p.Name, p.Host, p.Port, p.Username, p.Password, p.Env, p.Colour, p.Database)
	if err != nil {
		return false, errors.Wrap(err, "failed to insert new connection in postgres")
	}

	return true, nil
}

func (c *Connections) RefreshPostgresDatabase(id int64, dbID, dbName, poolID string) (*model.Database, error) {
	poolIDUUID, err := uuid.Parse(poolID)
	if err != nil {
		return nil, err
	}

	// Get all tables
	tables, err := c.GetAllPostgresTables(poolIDUUID)
	if err != nil {
		return nil, err
	}

	return &model.Database{
		ID:                   dbID,
		Name:                 dbName,
		PostgresConnectionID: id,
		PoolID:               poolID,
		IsActive:             true,
		Tables:               tables,
	}, nil
}

// This func is used to connect a specific database within a server
// id is the postgres connection id primary key in the sqlite3 database
// dbID uniquely identifies the active database within a connection
func (c *Connections) EstablishPostgresDatabaseConnection(id int64, dbName string) (*model.Database, error) {
	// Query for a single row by ID
	var name, host, port, username, password, colour string
	row := c.DB.QueryRow("SELECT name, host, port, username, password, colour FROM postgres WHERE id = ?", id)

	// Scan the result into variables
	err := row.Scan(&name, &host, &port, &username, &password, &colour)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(err, "postgres server not found")
		} else {
			return nil, err
		}
	}

	// Make a connection string
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require", username, password, host, port, dbName)

	activePoolID := uuid.New()

	// Establish connection and add pool to active pool manager
	_, err = c.PM.AddPool(activePoolID, connString)
	if err != nil {
		return nil, err
	}

	// Get all table names for suggestions
	tables, err := c.GetAllPostgresTables(activePoolID)
	if err != nil {
		return nil, err
	}

	// Get all columns of all tables for suggestions
	columns, err := c.GetAllDatabaseColumns(activePoolID)
	if err != nil {
		return nil, err
	}

	activeDB := name + " - " + dbName

	// Save the active db properties in all the tabs with type editor if active db properties are null
	_, err = c.DB.Exec("UPDATE tabs SET active_db_id = ?, active_db = ?, active_db_colour = ? WHERE active_db_id IS NULL AND type = 'editor'", activePoolID.String(), activeDB, colour)
	if err != nil {
		return nil, err
	}

	// In case of table rows, find all the rows with type table where
	// active_db_id is null and postgres_connection_id and database matches
	// set the active pool id and active db properties in such tabs
	_, err = c.DB.Exec("UPDATE tabs SET active_db_id = ?, active_db = ?, active_db_colour = ? WHERE active_db_id IS NULL AND type = 'table' AND postgres_conn_id = ? AND db_name = ?", activePoolID.String(), activeDB, colour, id, dbName)
	if err != nil {
		return nil, err
	}

	return &model.Database{
		Name:                   dbName,
		PostgresConnectionID:   id,
		PostgresConnectionName: name,
		Colour:                 colour,
		PoolID:                 activePoolID.String(),
		IsActive:               true,
		Tables:                 tables,
		Columns:                columns,
	}, nil
}

// This func is used to connect to a server
func (c *Connections) EstablishPostgresConnection(id int64) ([]model.Database, error) {
	// Query for a single row by ID
	var name, host, port, username, password, database, colour string
	row := c.DB.QueryRow("SELECT name, host, port, username, password, database, colour FROM postgres WHERE id = ?", id)

	// Scan the result into variables
	err := row.Scan(&name, &host, &port, &username, &password, &database, &colour)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(err, "postgres server not found")
		} else {
			return nil, err
		}
	}

	database = strings.TrimSpace(database)

	if database == "" {
		database = "postgres"
	}

	// Make a connection string
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require", username, password, host, port, database)

	activePoolID := uuid.New()

	// Establish connection and add pool to active pool manager
	_, err = c.PM.AddPool(activePoolID, connString)
	if err != nil {
		return nil, err
	}

	activeDB := name + " - " + database

	// Save the active db properties in all the tabs with type editor if active db properties are null
	_, err = c.DB.Exec("UPDATE tabs SET active_db_id = ?, active_db = ?, active_db_colour = ? WHERE active_db_id IS NULL AND type = 'editor'", activePoolID.String(), activeDB, colour)
	if err != nil {
		return nil, err
	}

	// In case of table rows, find all the rows with type table where
	// active_db_id is null and postgres_connection_id and database matches
	// set the active pool id and active db properties in such tabs
	_, err = c.DB.Exec("UPDATE tabs SET active_db_id = ?, active_db = ?, active_db_colour = ? WHERE active_db_id IS NULL AND type = 'table' AND postgres_conn_id = ? AND db_name = ?", activePoolID.String(), activeDB, colour, id, database)
	if err != nil {
		return nil, err
	}

	return c.GetPostgresServerDatabases(id, activePoolID, database, name, colour)
}

// Here, pool means the active connection to the database server
func (c *Connections) GetPostgresServerDatabases(postgresConnectionID int64, activePoolID uuid.UUID, activeDatabase, postgresConnectionName, colour string) ([]model.Database, error) {
	pool, exists := c.PM.GetPool(activePoolID)
	if !exists {
		return nil, errors.New("pool doesn't exist")
	}

	// Get Tables of active database
	tables, err := c.GetAllPostgresTables(activePoolID)
	if err != nil {
		return nil, err
	}

	// Get Columns of active database
	columns, err := c.GetAllDatabaseColumns(activePoolID)
	if err != nil {
		return nil, err
	}

	// Get all database names of the active connection
	rows, err := pool.Query(context.TODO(), "SELECT datname FROM pg_database")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Slice to hold results
	var databases []model.Database

	// Iterate through the rows
	for rows.Next() {
		var database model.Database
		err := rows.Scan(&database.Name)
		if err != nil {
			return nil, err
		}
		database.ID = "db_" + uuid.New().String()
		database.PostgresConnectionID = postgresConnectionID
		database.PostgresConnectionName = postgresConnectionName
		database.Colour = colour
		if database.Name == activeDatabase {
			database.PoolID = activePoolID.String()
			database.IsActive = true
			database.Tables = tables
			database.Columns = columns
		}

		databases = append(databases, database)
	}

	// Check for any error encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return databases, nil
}

func (c *Connections) GetAllPostgresTables(activePoolID uuid.UUID) ([]string, error) {
	pool, exists := c.PM.GetPool(activePoolID)
	if !exists {
		return nil, errors.New("pool doesn't exist")
	}

	// Get all tables
	rows, err := pool.Query(context.TODO(), "SELECT tablename FROM pg_tables WHERE schemaname = 'public'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Slice to hold results
	var tables []string

	// Iterate through the rows
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	// Check for any error encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

// Get all columns of the active database across all tables
func (c *Connections) GetAllDatabaseColumns(activePoolID uuid.UUID) ([]string, error) {
	pool, exists := c.PM.GetPool(activePoolID)
	if !exists {
		return nil, errors.New("pool doesn't exist")
	}

	// Get all columns
	query := `
		SELECT DISTINCT
			column_name
		FROM
			information_schema.columns
		WHERE
			table_schema NOT IN ('information_schema', 'pg_catalog')
	`
	rows, err := pool.Query(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Slice to hold results
	var columns []string

	// Iterate through the rows
	for rows.Next() {
		var column string
		err := rows.Scan(&column)
		if err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}

	// Check for any error encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return columns, nil
}

func (c *Connections) TerminatePostgresDatabaseConnection(activePoolID string) (bool, error) {
	activePoolIDUUID, err := uuid.Parse(activePoolID)
	if err != nil {
		return false, err
	}
	// Remove the db pool from active pools
	err = c.PM.DeletePool(activePoolIDUUID)
	if err != nil {
		return false, err
	}

	// Remove the pool from all the tabs in which it's saved
	_, err = c.DB.Exec("UPDATE tabs SET active_db_id = NULL, active_db = NULL, active_db_colour = NULL WHERE active_db_id = ?", activePoolID)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Connections) TerminateAllDatabaseConnections() error {
	c.PM.mu.Lock()
	defer c.PM.mu.Unlock()

	activeDBIds := []string{}

	for id, pool := range c.PM.Pools {
		activeDBIds = append(activeDBIds, id.String())
		pool.Close()
		delete(c.PM.Pools, id)
	}

	// Build placeholders (?, ?, ?)
	placeholders := strings.Repeat("?,", len(activeDBIds))
	placeholders = strings.TrimRight(placeholders, ",")

	// Remove the pool from all the tabs in which it's saved
	// Construct the query
	query := fmt.Sprintf(
		"UPDATE tabs SET active_db_id = NULL, active_db = NULL, active_db_colour = NULL WHERE active_db_id IN (%s)",
		placeholders,
	)

	// Convert []string to []interface{}
	args := make([]interface{}, len(activeDBIds))
	for i, v := range activeDBIds {
		args[i] = v
	}

	// Execute the query with placeholders
	_, err := c.DB.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

// Function to check if a query is a write operation
func isWriteOperation(query string) bool {
	// List of SQL keywords for write operations in lowercase
	writeKeywords := []string{"insert", "update", "delete", "alter", "create", "drop", "truncate"}

	// Convert query to lowercase for case-insensitive comparison
	query = strings.ToLower(strings.TrimSpace(query))

	// Check if the query starts with any write keyword
	for _, keyword := range writeKeywords {
		if strings.HasPrefix(query, keyword) {
			return true
		}
	}
	return false
}

func (c *Connections) ExecuteQuery(activePoolID uuid.UUID, query string, tabID int64) *model.QueryResult {
	pool, exists := c.PM.GetPool(activePoolID)
	if !exists {
		return &model.QueryResult{OK: false, Message: "pool doesn't exist"}
	}

	ctx := context.Background()

	response := &model.QueryResult{OK: true}

	normalizedQuery := strings.ToLower(strings.TrimSpace(query))

	isWrite := isWriteOperation(normalizedQuery)

	if isWrite {
		// Use Exec for write operations
		tag, err := pool.Exec(ctx, query)
		if err != nil {
			return &model.QueryResult{
				OK:           true,
				Message:      err.Error(),
				RowsAffected: int64(0),
				Columns:      []string{"Error"},
				Rows:         [][]model.Cell{{model.Cell{Column: "Error", Value: err.Error()}}},
			}
		}
		response.RowsAffected = tag.RowsAffected()
		response.Columns = []string{"Rows Affected"}
		response.Rows = [][]model.Cell{{model.Cell{Column: "Rows Affected", Value: fmt.Sprintf("%d", response.RowsAffected)}}}
	} else {
		// Use Query for read operations
		resultRows, err := pool.Query(ctx, query)
		if err != nil {
			return &model.QueryResult{
				OK:           true,
				Message:      err.Error(),
				RowsAffected: int64(0),
				Columns:      []string{"Error"},
				Rows:         [][]model.Cell{{model.Cell{Column: "Error", Value: err.Error()}}},
			}
		}
		defer resultRows.Close()

		columns := resultRows.FieldDescriptions()
		columnNames := make([]string, len(columns))
		for i, column := range columns {
			columnNames[i] = string(column.Name)
		}
		response.Columns = columnNames

		var rows [][]model.Cell

		for resultRows.Next() {
			row, err := resultRows.Values()
			if err != nil {
				return &model.QueryResult{OK: false, Message: err.Error()}
			}

			cells := []model.Cell{}
			for i, cell := range row {
				newCell := model.Cell{
					Column: columnNames[i],
				}
				switch v := cell.(type) {
				case []byte:
					newCell.Value = string(v)
				case time.Time:
					newCell.Value = v.Format(time.RFC3339)
				case nil:
					newCell.Value = "NULL"
				case [16]uint8:
					newCell.Value = uuid.UUID(v).String()
				case string:
					if v == "" {
						newCell.Value = "EMPTY"
					}
					newCell.Value = v
				default:
					newCell.Value = fmt.Sprintf("%v", v)
				}
				cells = append(cells, newCell)
			}
			rows = append(rows, cells)
		}

		if err := resultRows.Err(); err != nil {
			return &model.QueryResult{OK: false, Message: err.Error()}
		}

		response.Rows = rows
	}

	output := &model.Output{
		Columns: response.Columns,
		Rows:    response.Rows,
	}

	go c.UpdateTabOutput(tabID, output)

	return response
}

func (c *Connections) GetTableData(activePoolID uuid.UUID, tabID int64, tableName, selectQuery, limit, offset, where, orderBy, groupBy string) *model.QueryResult {
	pool, exists := c.PM.GetPool(activePoolID)
	if !exists {
		return &model.QueryResult{OK: false, Message: "pool doesn't exist"}
	}

	ctx := context.Background()

	response := &model.QueryResult{OK: true}

	setLimit := strconv.Itoa(100)
	if strings.TrimSpace(limit) != "" {
		limitInt, err := strconv.Atoi(strings.TrimSpace(limit))
		if err != nil {
			return &model.QueryResult{OK: false, Message: "limit is not a number"}
		}
		if limitInt > 100 {
			return &model.QueryResult{OK: false, Message: "limit cannot be greater than 100"}
		}
		setLimit = strings.TrimSpace(limit)
	}

	if strings.TrimSpace(selectQuery) == "" {
		selectQuery = "*"
	}

	query := fmt.Sprintf("SELECT %s FROM %s", selectQuery, tableName)
	if strings.TrimSpace(where) != "" {
		query += fmt.Sprintf(" WHERE %s", strings.TrimSpace(where))
	}
	if strings.TrimSpace(groupBy) != "" {
		query += fmt.Sprintf(" GROUP BY %s", strings.TrimSpace(groupBy))
	}
	if strings.TrimSpace(orderBy) != "" {
		query += fmt.Sprintf(" ORDER BY %s", strings.TrimSpace(orderBy))
	}

	// Set limit
	query += fmt.Sprintf(" LIMIT %s", setLimit)

	if strings.TrimSpace(offset) != "" {
		query += fmt.Sprintf(" OFFSET %s", strings.TrimSpace(offset))
	}

	// Use Query for read operations
	resultRows, err := pool.Query(ctx, query)
	if err != nil {
		return &model.QueryResult{
			OK:           true,
			Message:      err.Error(),
			RowsAffected: int64(0),
			Columns:      []string{"Error"},
			Rows:         [][]model.Cell{{model.Cell{Column: "Error", Value: err.Error()}}},
		}
	}
	defer resultRows.Close()

	columns := resultRows.FieldDescriptions()
	columnNames := make([]string, len(columns))
	for i, column := range columns {
		columnNames[i] = string(column.Name)
	}
	response.Columns = columnNames

	var rows [][]model.Cell

	for resultRows.Next() {
		row, err := resultRows.Values()
		if err != nil {
			return &model.QueryResult{OK: false, Message: err.Error()}
		}

		cells := []model.Cell{}
		for i, cell := range row {
			newCell := model.Cell{
				Column: columnNames[i],
			}
			switch v := cell.(type) {
			case []byte:
				newCell.Value = string(v)
			case time.Time:
				newCell.Value = v.Format(time.RFC3339)
			case nil:
				newCell.Value = "NULL"
			case [16]uint8:
				newCell.Value = uuid.UUID(v).String()
			case string:
				if v == "" {
					newCell.Value = "EMPTY"
				}
				newCell.Value = v
			default:
				newCell.Value = fmt.Sprintf("%v", v)
			}
			cells = append(cells, newCell)
		}
		rows = append(rows, cells)
	}

	if err := resultRows.Err(); err != nil {
		return &model.QueryResult{OK: false, Message: err.Error()}
	}

	response.Rows = rows

	output := &model.Output{
		Columns: response.Columns,
		Rows:    response.Rows,
	}

	go c.UpdateTabOutput(tabID, output)

	return response
}

func (c *Connections) UpdateTabOutput(tabID int64, output *model.Output) {
	jsonOutput, err := json.Marshal(output)
	if err != nil {
		fmt.Println(err)
	}

	query := `UPDATE tabs SET output = ? WHERE id = ?`
	_, err = c.DB.Exec(query, string(jsonOutput), tabID)
	if err != nil {
		fmt.Println(err)
	}
}
