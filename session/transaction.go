package session

import "tsorm/log"

// Begin starts a transaction.
func (s *Session) Begin() (err error) {
	// Log the beginning of the transaction.
	log.Info("transaction begin")
	// Start the transaction.
	if s.tx, err = s.db.Begin(); err != nil {
		log.Error(err)
		return
	}
	return
}

// Commit commits the transaction.
func (s *Session) Commit() (err error) {
	// Log the transaction commit.
	log.Info("transaction commit")
	// Commit the transaction.
	if err = s.tx.Commit(); err != nil {
		log.Error(err)
	}
	return
}

// Rollback rolls back the transaction.
func (s *Session) Rollback() (err error) {
	// Log the transaction rollback.
	log.Info("transaction rollback")
	// Roll back the transaction.
	if err = s.tx.Rollback(); err != nil {
		log.Error(err)
	}
	return
}
