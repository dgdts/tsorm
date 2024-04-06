package session

import (
	"testing"
	"tsorm/log"
)

// Account represents a user account.
type Account struct {
	ID       int    `tsorm:"PRIMARY KEY"` // ID is the primary key of the account.
	Password string // Password is the password of the account.
}

// BeforeInsert is a callback method that is executed before inserting a record into the database.
// It increments the ID of the account by 1000 before insertion.
func (account *Account) BeforeInsert(s *Session) error {
	log.Info("before insert", account)
	account.ID += 1000
	return nil
}

// AfterQuery is a callback method that is executed after querying a record from the database.
// It obfuscates the password of the account by replacing it with asterisks.
func (account *Account) AfterQuery(s *Session) error {
	log.Info("after query", account)
	account.Password = "******"
	return nil
}

// TestSession_CallMethod tests the CallMethod function of the Session struct.
func TestSession_CallMethod(t *testing.T) {
	// Create a new session for testing.
	s := NewSessionForTest(t).Model(&Account{})

	// Drop the table if it exists.
	_ = s.DropTable()
	// Create the table.
	_ = s.CreateTable()
	// Insert test data into the table.
	_, _ = s.Insert(&Account{ID: 0, Password: "456"}, &Account{ID: 1, Password: "321"})

	// Create an empty account object.
	u := &Account{}

	// Query the first record into the account object.
	err := s.First(u)
	// Check for errors or unexpected values.
	if err != nil || u.ID != 1000 || u.Password != "******" {
		t.Fatal("test failed, got", u)
	}
}
