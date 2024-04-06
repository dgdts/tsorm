package dialect

import (
	"fmt"
	"reflect"
	"time"
)

// sqlite3 represents the SQLite3 dialect.
type sqlite3 struct{}

// Ensure that sqlite3 implements the Dialect interface.
var _ Dialect = (*sqlite3)(nil)

// DataTypeOf returns the corresponding SQL data type for the given Go type.
func (s *sqlite3) DataTypeOf(t reflect.Value) string {
	switch t.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		// Check if the struct type is time.Time, if so, return "datetime".
		if _, ok := t.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	// Panic if the type is not supported.
	panic(fmt.Sprintf("invalid SQL type %s (%s)", t.Type().Name(), t.Kind()))
}

// TableExistSQL returns the SQL query to check if a table exists, along with any associated variables.
func (s *sqlite3) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	// SQL query to check if the table exists in the SQLite master table.
	return "SELECT name FROM sqlite_master WHERE type='table' and name = ?", args
}

// init registers the sqlite3 dialect when the package is initialized.
func init() {
	RegisterDialect("sqlite3", &sqlite3{})
}
