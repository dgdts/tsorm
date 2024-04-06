package clause

import (
	"fmt"
	"strings"
)

// generator defines the function type for generating SQL statements.
type generator func(values ...interface{}) (string, []interface{})

// generators is a map that associates SQL types (Type) with their respective generator functions.
var generators map[Type]generator

// init function is automatically executed when the package is imported, used to initialize the generators map.
func init() {
	generators = make(map[Type]generator)
	// Associates various generator functions with corresponding SQL types.
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderby
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

// genBindVars generates the binding variable string, where 'num' specifies the number of binding variables.
func genBindVars(num int) string {
	var vars []string
	// Generates the specified number of binding variable strings "?" in a loop.
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	// Joins the binding variables with commas to form a string and returns.
	return strings.Join(vars, ", ")
}

// _insert generates the SQL string and related variables for the INSERT statement.
func _insert(values ...interface{}) (string, []interface{}) {
	// Parses the input parameters, where tableName represents the table name and fields represents the field names.
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	// Returns the formatted INSERT statement and an empty variable slice, as variables in the VALUES clause are handled in the _values function.
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []interface{}{}
}

// _values generates the SQL string and related variables for the VALUES clause.
func _values(values ...interface{}) (string, []interface{}) {
	var bingStr string         // Stores the binding variable string
	var sql strings.Builder    // Used to build the SQL string
	var vars []interface{}     // Stores the related variables
	sql.WriteString("VALUES ") // Appends the VALUES clause
	// Iterates over the input values
	for i, value := range values {
		v := value.([]interface{}) // Converts the parameter to an interface slice
		// If bingStr is empty, generates a binding variable string with the same number of values
		if bingStr == "" {
			bingStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bingStr)) // Appends the VALUES clause
		if i+1 != len(values) {
			sql.WriteString(", ") // Adds a comma if it's not the last value
		}
		vars = append(vars, v...) // Adds the values to the variable slice
	}
	// Returns the generated SQL string and related variable slice
	return sql.String(), vars
}

// _select generates the SQL string and related variables for the SELECT statement.
func _select(values ...interface{}) (string, []interface{}) {
	// Parses the input parameters, where tableName represents the table name and fields represents the field names.
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	// Returns the formatted SELECT statement and an empty variable slice, as there are no related variables.
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), []interface{}{}
}

// _limit generates the SQL string and related variables for the LIMIT clause.
func _limit(values ...interface{}) (string, []interface{}) {
	// Returns the formatted LIMIT clause and the provided variables, without further processing.
	return "LIMIT ?", values
}

// _where generates the SQL string and related variables for the WHERE clause.
func _where(values ...interface{}) (string, []interface{}) {
	// Parses the input parameters, where desc represents the WHERE condition description and vars represents the variables in the WHERE condition.
	desc, vars := values[0], values[1:]
	// Returns the formatted WHERE clause and the related variable slice.
	return fmt.Sprintf("WHERE %s", desc), vars
}

// _orderby generates the SQL string and related variables for the ORDER BY clause.
func _orderby(values ...interface{}) (string, []interface{}) {
	// Returns the formatted ORDER BY clause and an empty variable slice, as there are no related variables.
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}

// _update generates the SQL string and related variables for the UPDATE statement.
func _update(values ...interface{}) (string, []interface{}) {
	// Parses the input parameters, where tableName represents the table name and fieldNames represents the field names and their corresponding values.
	tableName := values[0]
	fieldNames := values[1].(map[string]interface{})
	var keys []string      // Stores the strings of field names and values
	var vars []interface{} // Stores the related variables
	// Iterates over the field names and values, building the SET clause.
	for k, v := range fieldNames {
		keys = append(keys, k+" = ?")
		vars = append(vars, v) // Adds the value to the variable slice
	}
	// Returns the formatted UPDATE statement and related variable slice.
	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ", ")), vars
}

// _delete generates the SQL string and related variables for the DELETE statement.
func _delete(values ...interface{}) (string, []interface{}) {
	// Returns the formatted DELETE statement and provided variables, without further processing.
	return fmt.Sprintf("DELETE FROM %s", values[0]), []interface{}{}
}

// _count generates the SQL string and related variables for the COUNT function.
func _count(values ...interface{}) (string, []interface{}) {
	// Calls the _select function to generate the SELECT COUNT(*) statement and returns.
	return _select(values[0], []string{"COUNT(*)"})
}
