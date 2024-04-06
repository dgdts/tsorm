package tsorm

import (
	"errors"
	"reflect"
	"testing"
	"tsorm/session"

	_ "github.com/mattn/go-sqlite3"
)

// OpenDB opens a database connection for testing purposes.
func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("sqlite3", "ts.db")
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	return engine
}

// User represents a user entity in the database.
type User struct {
	Name string `tsorm:"PRIMARY KEY"`
	Age  int
}

// TestEngine_TransactionRollBack tests transaction rollback functionality.
func TestEngine_TransactionRollBack(t *testing.T) {
	engline := OpenDB(t)
	defer engline.Close()
	s := engline.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := engline.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		s.Insert(&User{Name: "Tom", Age: 10})
		return nil, errors.New("Error")
	})
	if err == nil || s.HasTable() {
		t.Fatal("failed to rollback")
	}
}

// TestEngine_TransactionCommit tests transaction commit functionality.
func TestEngine_TransactionCommit(t *testing.T) {
	engline := OpenDB(t)
	defer engline.Close()
	s := engline.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := engline.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		ret, err := s.Insert(&User{Name: "Tom", Age: 10})
		return ret, err
	})

	u := &User{}
	s.First(u)

	if err != nil || u.Name != "Tom" {
		t.Fatal("failed to commit")
	}
}

// TestEngine_Migrate tests schema migration functionality.
func TestEngine_Migrate(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()

	s := engine.NewSession()
	s.Raw("DROP TABLE IF EXISTS User;").Exec()
	s.Raw("CREATE TABLE User(Name text PRIMARY KEY, XXX integer);").Exec()
	s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	engine.Migrate(&User{})

	rows, _ := s.Raw("SELECT * FROM User").QueryRows()
	columns, _ := rows.Columns()
	if !reflect.DeepEqual(columns, []string{"Name", "Age"}) {
		t.Fatal("Failed to migrate table User, got columns", columns)
	}
}
