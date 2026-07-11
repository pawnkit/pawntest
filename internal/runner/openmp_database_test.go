package runner

import (
	"testing"

	"github.com/pawnkit/pawntest/internal/backend"
)

func TestDatabaseScenarioExecutesSQLiteQueries(t *testing.T) {
	vm, _ := registeredScenarios(t)
	vm.strings[100] = ":memory:"
	database := callScenarioNative(t, vm, "DB_Open", 100)
	if database == 0 {
		t.Fatal("DB_Open failed")
	}

	queries := []string{
		"CREATE TABLE players (id INTEGER, name TEXT, score REAL)",
		"INSERT INTO players VALUES (1, 'Alice', 9.5), (2, 'Bob', 7.25)",
	}
	for index, query := range queries {
		address := backend.Cell(200 + index)
		vm.strings[address] = query
		result := callScenarioNative(t, vm, "DB_ExecuteQuery", database, address)
		if result == 0 || callScenarioNative(t, vm, "DB_FreeResultSet", result) != 1 {
			t.Fatalf("query %d failed", index)
		}
	}

	vm.strings[300] = "SELECT id, name, score FROM players ORDER BY id"
	result := callScenarioNative(t, vm, "DB_ExecuteQuery", database, 300)
	if rows := callScenarioNative(t, vm, "DB_GetRowCount", result); rows != 2 {
		t.Fatalf("DB_GetRowCount = %d, want 2", rows)
	}
	if fields := callScenarioNative(t, vm, "DB_GetFieldCount", result); fields != 3 {
		t.Fatalf("DB_GetFieldCount = %d, want 3", fields)
	}
	if id := callScenarioNative(t, vm, "DB_GetFieldInt", result, 0); id != 1 {
		t.Fatalf("first id = %d, want 1", id)
	}
	callScenarioNative(t, vm, "DB_GetFieldString", result, 1, 400, 16)
	if vm.strings[400] != "Alice" {
		t.Fatalf("first name = %q, want Alice", vm.strings[400])
	}
	vm.strings[500] = "name"
	if callScenarioNative(t, vm, "DB_SelectNextRow", result) != 1 {
		t.Fatal("DB_SelectNextRow failed")
	}
	callScenarioNative(t, vm, "DB_GetFieldStringByName", result, 500, 600, 16)
	if vm.strings[600] != "Bob" {
		t.Fatalf("second name = %q, want Bob", vm.strings[600])
	}

	if count := callScenarioNative(t, vm, "DB_GetDatabaseResultSetCount"); count != 1 {
		t.Fatalf("result count = %d, want 1", count)
	}
	callScenarioNative(t, vm, "DB_FreeResultSet", result)
	callScenarioNative(t, vm, "DB_Close", database)
	if count := callScenarioNative(t, vm, "DB_GetDatabaseConnectionCount"); count != 0 {
		t.Fatalf("connection count = %d, want 0", count)
	}
}

func TestDatabaseScenarioCloseReleasesResources(t *testing.T) {
	vm, registry := registeredScenarios(t)
	vm.strings[100] = ":memory:"
	if database := callScenarioNative(t, vm, "DB_Open", 100); database == 0 {
		t.Fatal("DB_Open failed")
	}

	state, ok := registry.modules[15].(*databaseState)
	if !ok {
		t.Fatal("scenario registry did not contain database state")
	}
	if err := state.Close(); err != nil {
		t.Fatal(err)
	}
	if len(state.connections) != 0 || len(state.results) != 0 {
		t.Fatal("database resources were not cleared")
	}
}
