package schema

import (
	"reflect"
	"testing"
	"tsorm/dialect"
)

// Define a sample struct for testing.
type User struct {
	ID   int    `tsorm:"PRIMARY KEY"`
	Name string `tsorm:"NOT NULL"`
	Age  int
}

var TestDial, _ = dialect.GetDialect("sqlite3")

// TestParse tests the Parse function of the schema package.
func TestParse(t *testing.T) {

	// Parse the schema for the User struct.
	s := Parse(&User{}, TestDial)

	// Check the schema's name.
	if s.Name != "User" {
		t.Errorf("Expected schema name 'User', got '%s'", s.Name)
	}

	// Check the number of fields in the schema.
	expectedNumFields := 3 // ID, Name, Age
	if len(s.Fields) != expectedNumFields {
		t.Errorf("Expected %d fields, got %d", expectedNumFields, len(s.Fields))
	}

	// Check the field types and tags.
	expectedFields := map[string]struct {
		Type string
		Tag  string
	}{
		"ID":   {"integer", "PRIMARY KEY"},
		"Name": {"text", "NOT NULL"},
		"Age":  {"integer", ""},
	}

	for _, field := range s.Fields {
		expectedField, ok := expectedFields[field.Name]
		if !ok {
			t.Errorf("Unexpected field: %s", field.Name)
		}
		if field.Type != expectedField.Type {
			t.Errorf("Expected type %s for field %s, got %s", expectedField.Type, field.Name, field.Type)
		}
		if field.Tag != expectedField.Tag {
			t.Errorf("Expected tag %s for field %s, got %s", expectedField.Tag, field.Name, field.Tag)
		}
	}
}

// TestRecordValues tests the RecordValues method of the Schema struct.
func TestRecordValues(t *testing.T) {
	// Create a sample user record.
	user := User{ID: 1, Name: "John", Age: 30}

	// Parse the schema for the User struct.
	s := Parse(&User{}, TestDial)

	// Get the record values from the user record.
	recordValues := s.RecordValues(user)

	// Check the number of record values.
	expectedNumValues := 3 // ID, Name, Age
	if len(recordValues) != expectedNumValues {
		t.Errorf("Expected %d record values, got %d", expectedNumValues, len(recordValues))
	}

	// Check the record values.
	expectedValues := []interface{}{1, "John", 30}
	for i, val := range recordValues {
		expectedVal := expectedValues[i]
		if !reflect.DeepEqual(val, expectedVal) {
			t.Errorf("Expected record value %v, got %v", expectedVal, val)
		}
	}
}
