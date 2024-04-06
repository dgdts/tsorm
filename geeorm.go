package tsorm

import (
	"database/sql"
	"fmt"
	"strings"
	"tsorm/dialect"
	"tsorm/log"
	"tsorm/session"
)

// Engine represents the database engine.
type Engine struct {
	db      *sql.DB         // Underlying database connection
	dialect dialect.Dialect // Database dialect
}

// NewEngine creates a new database engine.
func NewEngine(driver string, source string) (e *Engine, err error) {
	// Open a database connection.
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// Ping the database to ensure connectivity.
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	// Get the dialect for the specified driver.
	dial, ok := dialect.GetDialect("sqlite3")
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}

	// Initialize the Engine with the database connection and dialect.
	e = &Engine{
		db:      db,
		dialect: dial,
	}

	// Log successful database connection.
	log.Info("Connect database success")
	return
}

// Close closes the database engine.
func (e *Engine) Close() {
	// Close the underlying database connection.
	if err := e.db.Close(); err != nil {
		log.Error("Failed to close database:", err)
	} else {
		log.Info("Close database success")
	}
}

// NewSession creates a new session associated with the engine.
func (e *Engine) NewSession() *session.Session {
	return session.NewSession(e.db, e.dialect)
}

// TxFunc represents a function signature for transactions.
type TxFunc func(s *session.Session) (result interface{}, err error)

// Transaction performs a transaction with the provided function.
func (e *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := e.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p)
		} else if err != nil {
			_ = s.Rollback()
		} else {
			err = s.Commit()
		}
	}()
	return f(s)
}

// difference returns the difference between two string slices.
func difference(a []string, b []string) (diff []string) {
	mapB := make(map[string]bool)
	for _, v := range b {
		mapB[v] = true
	}

	for _, v := range a {
		if _, ok := mapB[v]; !ok {
			diff = append(diff, v)
		}
	}
	return
}

// Migrate migrates the schema of the given value to the database.
func (e *Engine) Migrate(value interface{}) error {
	_, err := e.Transaction(func(s *session.Session) (result interface{}, err error) {
		// If the table does not exist, create it.
		if !s.Model(value).HasTable() {
			log.Infof("table %s doesn't exist", s.RefTable().Name)
			return nil, s.CreateTable()
		}

		// Get the table schema.
		table := s.RefTable()

		// Query a row from the table to get the columns.
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
		columns, _ := rows.Columns()

		// Find columns to add and delete.
		addCols := difference(table.FieldNames, columns)
		delCols := difference(columns, table.FieldNames)
		log.Infof("added cols %v, deleted cols %v", addCols, delCols)

		// Add new columns.
		for _, col := range addCols {
			f := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.Name, f.Name, f.Type)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}

		// If no columns are to be deleted, return.
		if len(delCols) == 0 {
			return
		}

		// Rename the table and recreate it to delete columns.
		temp := "temp_" + table.Name
		fieldStr := strings.Join(table.FieldNames, ",")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s from %s;", temp, fieldStr, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", temp, table.Name))
		_, err = s.Exec()
		return
	})
	return err
}
