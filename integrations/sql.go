package integrations

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/logger"
	"github.com/SkySingh04/fractal/registry"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// PostgreSQLSource struct represents the configuration for consuming messages from PostgreSQL.
type PostgreSQLSource struct {
	ConnString string `json:"postgresql_source_conn_string"`
}

// PostgreSQLDestination struct represents the configuration for publishing messages to PostgreSQL.
type PostgreSQLDestination struct {
	ConnString string `json:"postgresql_target_conn_string"`
}

// FetchData connects to PostgreSQL, retrieves data, and returns it.
func (p PostgreSQLSource) FetchData(req interfaces.Request) (interface{}, error) {
	if req.SQLSourceConnString == "" {
		return nil, errors.New("missing PostgreSQL source connection string")
	}
	logger.Infof("Connecting to PostgreSQL source...")

	db, err := sql.Open("postgres", req.SQLSourceConnString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Retrieve the list of all tables in the public schema
	tablesQuery := "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'"
	rows, err := db.Query(tablesQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map to hold results categorized by table name
	allResults := make(map[string][]map[string]interface{})

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}

		// For each table, fetch its data
		dataQuery := "SELECT * FROM " + tableName // Fetch all columns from the table
		dataRows, err := db.Query(dataQuery)
		if err != nil {
			logger.Errorf("Error querying table %s: %s", tableName, err)
			continue // Skip this table on error
		}
		defer dataRows.Close()

		// Get column names for later use
		columns, err := dataRows.Columns()
		if err != nil {
			return nil, err
		}

		for dataRows.Next() {
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			if err := dataRows.Scan(valuePtrs...); err != nil {
				return nil, err
			}

			rowData := make(map[string]interface{})
			for i, colName := range columns {
				val := values[i]
				rowData[colName] = val
			}
			allResults[tableName] = append(allResults[tableName], rowData) // Append row data to the appropriate table key
		}

		if err := dataRows.Err(); err != nil {
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	logger.Infof("Data fetched from PostgreSQL: %v", allResults)
	return allResults, nil
}

// EnsureTableExistsWorker processes table creation tasks.
func EnsureTableExistsWorker(db *sql.DB, tasks chan map[string]interface{}, errorsChan chan error, done chan bool) {
	for task := range tasks {
		tableName := task["tableName"].(string)
		row := task["row"].(map[string]interface{})

		// Check if table exists
		checkQuery := fmt.Sprintf("SELECT to_regclass('public.%s')", tableName)
		var tableExists sql.NullString
		err := db.QueryRow(checkQuery).Scan(&tableExists)
		if err != nil {
			errorsChan <- err
			continue
		}

		// If table does not exist, create it
		if !tableExists.Valid {
			var columns []string
			for colName, value := range row {
				colType := "TEXT" // Default to TEXT type
				switch value.(type) {
				case int, int32, int64:
					colType = "INTEGER"
				case float32, float64:
					colType = "FLOAT"
				case bool:
					colType = "BOOLEAN"
				}
				columns = append(columns, fmt.Sprintf("%s %s", colName, colType))
			}
			createQuery := fmt.Sprintf("CREATE TABLE %s (%s)", tableName, strings.Join(columns, ", "))
			if _, err := db.Exec(createQuery); err != nil {
				errorsChan <- err
				continue
			}
		}
		errorsChan <- nil // Indicate success
	}
	done <- true
}

// EnsureTableExists enqueues table creation tasks and processes them concurrently.
func EnsureTableExists(db *sql.DB, tableName string, row map[string]interface{}) error {
	// Buffered channels to queue tasks and capture errors
	tasks := make(chan map[string]interface{}, 1)
	errorsChan := make(chan error, 1)
	done := make(chan bool)

	// Start a worker goroutine
	go EnsureTableExistsWorker(db, tasks, errorsChan, done)

	// Enqueue the task
	tasks <- map[string]interface{}{
		"tableName": tableName,
		"row":       row,
	}
	close(tasks) // Signal no more tasks

	// Wait for the worker to finish and check for errors
	<-done
	close(errorsChan)

	// Collect any errors
	for err := range errorsChan {
		if err != nil {
			return err
		}
	}
	return nil
}

// SendData connects to PostgreSQL and publishes data to the specified table.
func (p PostgreSQLDestination) SendData(data interface{}, req interfaces.Request) error {
	if req.SQLTargetConnString == "" {
		return errors.New("missing PostgreSQL target connection string")
	}
	logger.Infof("Connecting to PostgreSQL destination...")

	db, err := sql.Open("postgres", req.SQLTargetConnString)
	if err != nil {
		return err
	}
	defer db.Close()

	// Assert that data is a map with table names as keys and slices of maps as values
	dataMap, ok := data.(map[string][]map[string]interface{})
	if !ok {
		return errors.New("data must be a map with table names as keys and slices of maps as values")
	}

	for tableName, rows := range dataMap {
		for _, row := range rows {
			// Ensure the table exists
			if err := EnsureTableExists(db, tableName, row); err != nil {
				return err
			}

			// Prepare column names and values for the insert query
			var columns []string
			var placeholders []string
			var values []interface{}

			for colName, value := range row {
				columns = append(columns, colName)
				placeholders = append(placeholders, "$"+strconv.Itoa(len(values)+1))
				values = append(values, value)
			}

			// Construct the INSERT query
			query := "INSERT INTO " + tableName + " (" + strings.Join(columns, ", ") + ") VALUES (" + strings.Join(placeholders, ", ") + ")"

			if _, err := db.Exec(query, values...); err != nil {
				logger.Errorf("Error inserting into table %s: %s", tableName, err)
				return err // Return on error
			}
		}
	}

	return nil
}

// Initialize the PostgreSQL integrations by registering them with the registry.
func init() {
	registry.RegisterSource("PostgreSQL", PostgreSQLSource{})
	registry.RegisterDestination("PostgreSQL", PostgreSQLDestination{})
}
