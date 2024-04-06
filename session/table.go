package session

import (
	"fmt"
	"reflect"
	"strings"
	"tsorm/log"
	"tsorm/schema"
)

// Model sets the model for the session.
func (s *Session) Model(value interface{}) *Session {
	// If the reference table is not set or the type of the provided value differs from the reference table's type,
	// parse the value and set the reference table.
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

// RefTable returns the reference table of the session.
func (s *Session) RefTable() *schema.Schema {
	// If the reference table is nil, log an error and return nil.
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

// CreateTable creates a table in the database based on the schema of the reference table.
func (s *Session) CreateTable() error {
	table := s.RefTable()
	var columns []string
	// Construct column definitions for the table.
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	// Execute the SQL command to create the table.
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()
	return err
}

// DropTable drops the table from the database.
func (s *Session) DropTable() error {
	// Execute the SQL command to drop the table if it exists.
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec()
	return err
}

// HasTable checks if the table exists in the database.
func (s *Session) HasTable() bool {
	// Get the SQL command and its values to check table existence.
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	// Query the database to check if the table exists.
	row := s.Raw(sql, values...).QueryRow()
	var temp string
	_ = row.Scan(&temp)
	// Return true if the scanned table name matches the reference table's name.
	return temp == s.RefTable().Name
}
