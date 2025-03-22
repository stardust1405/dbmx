package app

import (
	"context"
	"database/sql"
	"dbmx/model"
	"fmt"
	"log"
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

func NewConnections(db *sql.DB) *Connections {
	return &Connections{
		DB: db,
		PM: NewPoolManager(),
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
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", p.Username, p.Password, p.Host, p.Port, p.Database)

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

func (c *Connections) EstablishPostgresDatabaseConnection(id int64, dbID, dbName string) (*model.Database, error) {
	// Query for a single row by ID
	var host, port, username, password string
	row := c.DB.QueryRow("SELECT host, port, username, password FROM postgres WHERE id = ?", id)

	// Scan the result into variables
	err := row.Scan(&host, &port, &username, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(err, "postgres server not found")
		} else {
			return nil, err
		}
	}

	// Make a connection string
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, dbName)

	activePoolID := uuid.New()

	// Establish connection and add pool to active pool manager
	_, err = c.PM.AddPool(activePoolID, connString)
	if err != nil {
		return nil, err
	}

	// Get all tables
	tables, err := c.GetAllPostgresTables(activePoolID)
	if err != nil {
		return nil, err
	}

	return &model.Database{
		ID:                   dbID,
		Name:                 dbName,
		PostgresConnectionID: id,
		PoolID:               activePoolID.String(),
		IsActive:             true,
		Tables:               tables,
	}, nil
}

func (c *Connections) EstablishPostgresConnection(id int64) ([]model.Database, error) {
	// Query for a single row by ID
	var host, port, username, password, database string
	row := c.DB.QueryRow("SELECT host, port, username, password, database FROM postgres WHERE id = ?", id)

	// Scan the result into variables
	err := row.Scan(&host, &port, &username, &password, &database)
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
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)

	activePoolID := uuid.New()

	// Establish connection and add pool to active pool manager
	_, err = c.PM.AddPool(activePoolID, connString)
	if err != nil {
		return nil, err
	}

	return c.GetPostgresServerDatabases(id, activePoolID, database)
}

// Here, pool means the active connection to the database server
func (c *Connections) GetPostgresServerDatabases(postgresConnectionID int64, activePoolID uuid.UUID, activeDatabase string) ([]model.Database, error) {
	pool, exists := c.PM.GetPool(activePoolID)
	if !exists {
		return nil, errors.New("pool doesn't exist")
	}

	// Get Tables of active database
	tables, err := c.GetAllPostgresTables(activePoolID)
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
		if database.Name == activeDatabase {
			database.PoolID = activePoolID.String()
			database.IsActive = true
			database.Tables = tables
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

func (c *Connections) TerminatePostgresConnection(activePoolID uuid.UUID) (bool, error) {
	// Start a transaction
	tx, err := c.DB.Begin()
	if err != nil {
		return false, err
	}

	// Remove the db pool from active pools
	err = c.PM.DeletePool(activePoolID)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return false, errors.Wrap(err, "failed to commit tx")
	}

	return true, nil
}

func (c *Connections) ExecuteQuery(activePoolID uuid.UUID, query string, args ...interface{}) *model.GenericResponse {
	pool, exists := c.PM.GetPool(activePoolID)
	if !exists {
		return errorResponse(errors.New("pool doesn't exist"))
	}

	ctx := context.Background()
	response := &model.GenericResponse{OK: true}

	normalizedQuery := strings.ToLower(strings.TrimSpace(query))
	isReadOperation := strings.HasPrefix(normalizedQuery, "select") ||
		strings.HasPrefix(normalizedQuery, "with") ||
		strings.HasPrefix(normalizedQuery, "show")

	if isReadOperation {
		// Use Query for read operations
		rows, err := pool.Query(ctx, query, args...)
		if err != nil {
			return errorResponse(err)
		}
		defer rows.Close()

		columns := rows.FieldDescriptions()
		columnNames := make([]string, len(columns))
		for i, col := range columns {
			columnNames[i] = string(col.Name)
		}

		for rows.Next() {
			values, err := rows.Values()
			if err != nil {
				return errorResponse(err)
			}

			rowData := make(map[string]interface{})
			for i, val := range values {
				switch v := val.(type) {
				case []byte:
					rowData[columnNames[i]] = string(v)
				case time.Time:
					rowData[columnNames[i]] = v.Format(time.RFC3339)
				default:
					rowData[columnNames[i]] = val
				}
			}
			response.Data = append(response.Data, rowData)
		}

		if err := rows.Err(); err != nil {
			return errorResponse(err)
		}
	} else {
		// Use Exec for write operations
		tag, err := pool.Exec(ctx, query, args...)
		if err != nil {
			return errorResponse(err)
		}
		response.RowsAffected = tag.RowsAffected()
	}

	return response
}

func errorResponse(err error) *model.GenericResponse {
	return &model.GenericResponse{
		OK:      false,
		Message: err.Error(),
	}
}

// func (gr *model.GenericResponse) PrettyPrint() {
// 	if gr.OK {
// 		if len(gr.Data) > 0 {
// 			fmt.Println("Query Results:")
// 			for i, row := range gr.Data {
// 				jsonData, _ := json.MarshalIndent(row, "  ", "  ")
// 				fmt.Printf("Row %d:\n%s\n", i+1, jsonData)
// 			}
// 		} else if gr.RowsAffected > 0 {
// 			fmt.Printf("Success. Rows affected: %d\n", gr.RowsAffected)
// 		} else {
// 			fmt.Println("Operation completed (no results)")
// 		}
// 	} else {
// 		fmt.Printf("Error: %s\n", gr.Message)
// 	}
// }
