package schema

import (
	"go/ast"
	"reflect"
	"tsorm/dialect"
)

// Field represents a field in a database schema.
type Field struct {
	Name string // Name of the field
	Type string // Type of the field
	Tag  string // Tag of the field
}

// Schema represents the schema of a database table.
type Schema struct {
	Model      interface{}       // Model is the struct type used to define the schema
	Name       string            // Name is the name of the table
	Fields     []*Field          // Fields is a slice of fields in the table
	FieldNames []string          // FieldNames is a slice of field names in the table
	fieldMap   map[string]*Field // fieldMap is a map of field names to Field objects
}

// GetField returns the field with the given name from the schema.
func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}

// Parse parses the schema for the given model using the specified dialect.
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modeType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modeType.Name(),
		fieldMap: make(map[string]*Field),
	}

	// Iterate over the fields of the model type.
	for i := 0; i < modeType.NumField(); i++ {
		p := modeType.Field(i)
		// Check if the field is exported and not anonymous.
		if !p.Anonymous && ast.IsExported(p.Name) {
			// Create a Field object for the field.
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			// Check if the field has a "tsorm" tag.
			if v, ok := p.Tag.Lookup("tsorm"); ok {
				field.Tag = v
			}
			// Add the field to the schema's Fields and fieldMap.
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

// RecordValues extracts field values from a record and returns them as a slice of interfaces.
func (s *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	// Iterate over the fields of the schema.
	for _, field := range s.Fields {
		// Get the value of the field from the record.
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
