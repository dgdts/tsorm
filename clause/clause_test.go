package clause

import (
	"reflect"
	"testing"
)

func TestSelect(t *testing.T) {
	// Create a new instance of the Clause struct.
	var clause Clause

	// Set the SQL clauses in the desired order.
	clause.Set(LIMIT, 3)                      // Set the LIMIT clause with a limit of 3.
	clause.Set(SELECT, "User", []string{"*"}) // Set the SELECT clause for the "User" table selecting all columns.
	clause.Set(WHERE, "Name = ?", "Tom")      // Set the WHERE clause to filter records where the Name is "Tom".
	clause.Set(ORDERBY, "Age ASC")            // Set the ORDER BY clause to order records by Age in ascending order.

	// Build the SQL statement with the specified order of SQL clauses.
	sql, vars := clause.Build(SELECT, WHERE, ORDERBY, LIMIT)
	t.Log(sql, vars)

	// Check if the built SQL statement matches the expected SQL statement.
	if sql != "SELECT * FROM User WHERE Name = ? ORDER BY Age ASC LIMIT ?" {
		t.Fatal("failed to build SQL")
	}

	// Check if the built SQL variables match the expected SQL variables.
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 3}) {
		t.Fatal("failed to build SQLVars")
	}
}

func TestIsValidType(t *testing.T) {
	tests := []struct {
		name string
		t    Type
		want bool
	}{
		{"ValidType", INSERT, true},       // Test with a valid SQL type.
		{"InvalidType", Type(100), false}, // Test with an invalid SQL type.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check if the isValidType function returns the expected result.
			got := isValidType(tt.t)
			if got != tt.want {
				t.Errorf("isValidType(%d) = %v; want %v", tt.t, got, tt.want)
			}
		})
	}
}
