package dialect

import (
	"reflect"
	"testing"
	"time"
)

// TestSQLite3DataTypeOf tests the DataTypeOf method of the sqlite3 dialect.
func TestSQLite3DataTypeOf(t *testing.T) {
	// Create a new instance of the sqlite3 dialect.
	sqlite := &sqlite3{}

	// Test cases with different Go types.
	testCases := []struct {
		input    interface{}
		expected string
	}{
		{true, "bool"},
		{int(1), "integer"},
		{int8(1), "integer"},
		{int16(1), "integer"},
		{int32(1), "integer"},
		{int64(1), "bigint"},
		{uint(1), "integer"},
		{uint8(1), "integer"},
		{uint16(1), "integer"},
		{uint32(1), "integer"},
		{uint64(1), "bigint"},
		{float32(1.0), "real"},
		{float64(1.0), "real"},
		{"text", "text"},
		{[]byte{1, 2, 3}, "blob"},
		{time.Now(), "datetime"},
	}

	// Iterate over test cases.
	for _, tc := range testCases {
		t.Run(reflect.TypeOf(tc.input).Name(), func(t *testing.T) {
			// Call the DataTypeOf method with the test input.
			dataType := sqlite.DataTypeOf(reflect.ValueOf(tc.input))

			// Compare the result with the expected value.
			if dataType != tc.expected {
				t.Errorf("got %s, want %s", dataType, tc.expected)
			}
		})
	}
}

// TestSQLite3TableExistSQL tests the TableExistSQL method of the sqlite3 dialect.
func TestSQLite3TableExistSQL(t *testing.T) {
	// Create a new instance of the sqlite3 dialect.
	sqlite := &sqlite3{}

	// Test case with a table name.
	tableName := "users"

	// Call the TableExistSQL method with the table name.
	sql, vars := sqlite.TableExistSQL(tableName)

	// Define the expected SQL query and variables.
	expectedSQL := "SELECT name FROM sqlite_master WHERE type='table' and name = ?"
	expectedVars := []interface{}{"users"}

	// Compare the result with the expected values.
	if sql != expectedSQL {
		t.Errorf("got %s, want %s", sql, expectedSQL)
	}
	if !reflect.DeepEqual(vars, expectedVars) {
		t.Errorf("got %v, want %v", vars, expectedVars)
	}
}
