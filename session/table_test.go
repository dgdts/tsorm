package session

import (
	"database/sql"
	"testing"
	"tsorm/dialect"

	_ "github.com/mattn/go-sqlite3"
)

// User represents a user entity.
type User struct {
	Name string `tsrom:"PRIMARY KEY"`
	Age  int
}

var (
	TestDB      *sql.DB
	TestDial, _ = dialect.GetDialect("sqlite3")
)

// NewSessionForTest creates a new session for testing purposes.
func NewSessionForTest(t *testing.T) *Session {
	// Open a new SQLite database connection.
	TestDB, err := sql.Open("sqlite3", "../ts.db")
	if err != nil {
		// If an error occurs while creating the database connection, fail the test.
		t.Fatal("Create DB Failed:", err)
	}
	// Return a new session using the created database connection and SQLite dialect.
	return NewSession(TestDB, TestDial)
}

// TestSession_CreateDropHasTable tests the CreateTable method of the Session.
func TestSession_CreateDropHasTable(t *testing.T) {
	// Create a new session for testing and specify the model as User.
	s := NewSessionForTest(t).Model(&User{})
	// Drop the User table if it exists.
	err := s.DropTable()
	if err != nil {
		t.Fatal("Drop table failed,", err.Error())
	}
	// Create the User table.
	err = s.CreateTable()
	if err != nil {
		t.Fatal("Creat table failed,", err.Error())
	}
	// Check if the User table has been created successfully.
	if !s.HasTable() {
		t.Fatal("Failed to create table User")
	}
}
