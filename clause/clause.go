package clause

import (
	"fmt"
	"strings"
)

// Clause represents a component for constructing SQL statements.
type Clause struct {
	sql     map[Type]string        // Stores the SQL statements for each SQL type
	sqlVars map[Type][]interface{} // Stores the SQL variables for each SQL type
}

// Type represents the type of SQL statement.
type Type int

// Constants defining various types of SQL statements.
const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

// Set method is used to set the SQL statement and variables for a specific type.
// If a SQL statement of the same type already exists, it will be overwritten.
func (c *Clause) Set(name Type, vars ...interface{}) {
	// Check if the provided SQL type is valid.
	if !isValidType(name) {
		panic(fmt.Sprintf("invalid SQL type: %d", name))
	}

	// Initialize maps if they are nil.
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}

	// Generate SQL statement and variables.
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

// Build method is used to construct the SQL statement.
// It takes a series of SQL types as parameters and constructs the corresponding SQL statement according to the specified order.
// It returns the constructed SQL statement and its associated variables.
func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	// Clean up maps after the method call.
	defer func() {
		c.sql = nil
		c.sqlVars = nil
	}()

	// Check if the provided SQL types are valid.
	for _, order := range orders {
		if !isValidType(order) {
			panic(fmt.Sprintf("invalid SQL type: %d", order))
		}
	}

	var sqls []string
	var vars []interface{}
	// Construct SQL statement according to the specified order.
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}

// isValidType function is used to check if the provided SQL type is valid.
// It returns true if valid, false otherwise.
func isValidType(t Type) bool {
	return t >= INSERT && t <= COUNT
}
