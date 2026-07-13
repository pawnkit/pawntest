package runner

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/pawnkit/pawntest/internal/backend"
	_ "modernc.org/sqlite"
)

type databaseResult struct {
	columns []string
	rows    [][]string
	row     int
}

type databaseState struct {
	nextConnection int
	connections    map[int]databaseConnection
	nextResult     int
	results        map[int]*databaseResult
	open           databaseOpener
}

func newDatabaseState() *databaseState {
	return &databaseState{connections: map[int]databaseConnection{}, results: map[int]*databaseResult{}, open: openSQLiteDatabase}
}

func (state *databaseState) Clone() scenarioModule {
	clone := newDatabaseState()
	clone.open = state.open

	return clone
}

func (state *databaseState) Close() error {
	var closeErrors []error
	for id, database := range state.connections {
		closeErrors = append(closeErrors, database.Close())
		delete(state.connections, id)
	}
	clear(state.results)

	return errors.Join(closeErrors...)
}

func (state *databaseState) StrictFailures() []string {
	failures := []string{}
	if len(state.results) != 0 {
		failures = append(failures, fmt.Sprintf("unfreed database results: %d", len(state.results)))
	}
	if len(state.connections) != 0 {
		failures = append(failures, fmt.Sprintf("unclosed database connections: %d", len(state.connections)))
	}

	return failures
}

func (state *databaseState) Register(vm backend.VM, context *executionContext) error {
	return registerScenarioNatives(vm, state.natives(context.state), context.mocks, context.allowUnknown)
}

func (state *databaseState) natives(result *nativeState) map[string]backend.NativeFunc {
	natives := map[string]backend.NativeFunc{
		"__pt_database_connections": assertScenarioCount(result, "database connections", func() int { return len(state.connections) }),
		"__pt_database_results":     assertScenarioCount(result, "database results", func() int { return len(state.results) }),
		"DB_Open":                   state.openDatabase, "DB_Close": state.closeDatabase, "DB_ExecuteQuery": state.executeQuery,
		"DB_FreeResultSet": state.freeResult, "DB_GetRowCount": state.getRowCount, "DB_SelectNextRow": state.selectNextRow,
		"DB_GetFieldCount": state.getFieldCount, "DB_GetFieldName": state.getFieldName,
		"DB_GetFieldString": state.getFieldString, "DB_GetFieldInt": state.getFieldInt,
		"DB_GetFieldFloat": state.getFieldFloat, "DB_GetFieldStringByName": state.getFieldStringByName,
		"DB_GetFieldIntByName": state.getFieldIntByName, "DB_GetFieldFloatByName": state.getFieldFloatByName,
		"DB_GetMemHandle": state.getDatabaseHandle, "DB_GetLegacyDBResult": state.getResultHandle,
		"DB_GetDatabaseConnectionCount": state.getConnectionCount, "DB_GetDatabaseResultSetCount": state.getResultCount,
	}
	aliases := map[string]string{
		"db_open": "DB_Open", "db_close": "DB_Close", "db_query": "DB_ExecuteQuery", "db_free_result": "DB_FreeResultSet",
		"db_num_rows": "DB_GetRowCount", "db_next_row": "DB_SelectNextRow", "db_num_fields": "DB_GetFieldCount",
		"db_field_name": "DB_GetFieldName", "db_get_field": "DB_GetFieldString", "db_get_field_int": "DB_GetFieldInt",
		"db_get_field_float": "DB_GetFieldFloat", "db_get_field_assoc": "DB_GetFieldStringByName",
		"db_get_field_assoc_int": "DB_GetFieldIntByName", "db_get_field_assoc_float": "DB_GetFieldFloatByName",
		"db_get_mem_handle": "DB_GetMemHandle", "db_get_result_mem_handle": "DB_GetLegacyDBResult",
		"db_debug_openfiles": "DB_GetDatabaseConnectionCount", "db_debug_openresults": "DB_GetDatabaseResultSetCount",
	}
	for alias, name := range aliases {
		natives[alias] = natives[name]
	}

	return natives
}

func (state *databaseState) openDatabase(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 {
		return 0, nil
	}
	name, err := readNativeParams(ctx, params).String(0)
	if err != nil {
		return 0, err
	}
	database, err := state.open(name)
	if err != nil {
		return 0, nil
	}
	if err := database.Ping(); err != nil {
		_ = database.Close()

		return 0, nil
	}
	id := state.nextConnection + 1
	state.nextConnection++
	state.connections[id] = database

	return backend.Cell(id), nil
}

func (state *databaseState) closeDatabase(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 {
		return 0, nil
	}
	database, ok := state.connections[int(params[0])]
	if !ok {
		return 0, nil
	}
	if err := database.Close(); err != nil {
		return 0, err
	}
	delete(state.connections, int(params[0]))

	return 1, nil
}

func (state *databaseState) executeQuery(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) < 2 {
		return 0, nil
	}
	database, ok := state.connections[int(params[0])]
	if !ok {
		return 0, nil
	}
	query, err := readNativeParams(ctx, params).String(1)
	if err != nil {
		return 0, err
	}
	result, err := queryDatabase(database, query)
	if err != nil {
		return 0, nil
	}
	id := state.nextResult + 1
	state.nextResult++
	state.results[id] = result

	return backend.Cell(id), nil
}

func queryDatabase(database databaseConnection, query string) (*databaseResult, error) {
	rows, err := database.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	result := &databaseResult{columns: columns}
	for rows.Next() {
		values := make([]any, len(columns))
		pointers := make([]any, len(columns))
		for index := range values {
			pointers[index] = &values[index]
		}
		if err := rows.Scan(pointers...); err != nil {
			return nil, err
		}
		row := make([]string, len(columns))
		for index, value := range values {
			if bytes, ok := value.([]byte); ok {
				row[index] = string(bytes)
			} else if value != nil {
				row[index] = fmt.Sprint(value)
			}
		}
		result.rows = append(result.rows, row)
	}

	return result, rows.Err()
}

type databaseOpener func(name string) (databaseConnection, error)

type databaseConnection interface {
	Ping() error
	Close() error
	Query(query string) (databaseRows, error)
}

type databaseRows interface {
	Columns() ([]string, error)
	Next() bool
	Scan(dest ...any) error
	Err() error
	Close() error
}

type sqliteDatabase struct {
	database *sql.DB
}

func openSQLiteDatabase(name string) (databaseConnection, error) {
	database, err := sql.Open("sqlite", name)
	if err != nil {
		return nil, err
	}
	database.SetMaxOpenConns(1)

	return sqliteDatabase{database: database}, nil
}

func (database sqliteDatabase) Ping() error {
	return database.database.Ping()
}

func (database sqliteDatabase) Close() error {
	return database.database.Close()
}

func (database sqliteDatabase) Query(query string) (databaseRows, error) {
	return database.database.Query(query)
}

func parseDatabaseInt(value string) backend.Cell {
	parsed, _ := strconv.ParseInt(value, 10, 32)

	return backend.Cell(parsed)
}

func parseDatabaseFloat(value string) backend.Cell {
	parsed, _ := strconv.ParseFloat(value, 32)

	return floatCell(float32(parsed))
}
