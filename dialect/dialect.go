package dialect

import (
	"reflect"
)

// Dialect represents an interface for defining SQL dialects.
type Dialect interface {
	// DataTypeOf returns the SQL data type corresponding to the given Go type.
	DataTypeOf(t reflect.Value) string

	// TableExistSQL returns the SQL query to check if a table exists, along with any associated variables.
	TableExistSQL(tableName string) (string, []interface{})
}

// dialectsMap is a map that stores registered dialects.
var dialectsMap = map[string]Dialect{}

// RegisterDialect registers a dialect with the given name.
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

// GetDialect retrieves the dialect with the given name from the dialectsMap.
// It returns the dialect and a boolean indicating whether the dialect was found.
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
