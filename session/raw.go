package session

import (
	"database/sql"
	"strings"
	"tsorm/clause"
	"tsorm/dialect"
	"tsorm/log"
	"tsorm/schema"
)

// Session represents a database session.
type Session struct {
	db       *sql.DB         // db is the underlying SQL database connection.
	dialect  dialect.Dialect // dialect is the SQL dialect used by the session.
	tx       *sql.Tx         // tx is the SQL transaction associated with the session.
	refTable *schema.Schema  // refTable is the schema of the model associated with the session.
	clause   clause.Clause   // clause represents the SQL clauses used by the session.
	sql      strings.Builder // sql is the SQL query being constructed.
	sqlVars  []interface{}   // sqlVars contains the values to be used in the SQL query.
}

// CommonDB represents the common methods shared by both *sql.DB and *sql.Tx.
type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// Check if *sql.DB and *sql.Tx implement the CommonDB interface.
var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

// NewSession creates a new session with the given SQL database and dialect.
func NewSession(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
		tx:      nil,
	}
}

// Clear resets the session's state by clearing the SQL query and variables, and resetting the clause.
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

// DB returns the underlying SQL database connection or transaction.
// If a transaction is active, it returns the transaction; otherwise, it returns the database connection.
func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

// Raw appends raw SQL query and values to the session's SQL query.
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// Exec executes the SQL query built by the session and returns the result.
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)

	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// QueryRow executes the SQL query built by the session and returns a single row result.
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// QueryRows executes the SQL query built by the session and returns multiple row results.
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}
