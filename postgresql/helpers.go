package postgresql

import (
	"database/sql"
	"log"

	"github.com/hashicorp/errwrap"
)

// QueryAble is a DB connection (sql.DB/Tx)
type QueryAble interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// startTransaction starts a new DB transaction on the specified database.
// If the database is specified and different from the one configured in the provider,
// it will create a new connection pool if needed.
func startTransaction(client *Client, database string) (*sql.Tx, error) {
	if database != "" && database != client.databaseName {
		var err error
		client, err = client.config.NewClient(database)
		if err != nil {
			return nil, err
		}
	}
	db := client.DB()
	txn, err := db.Begin()
	if err != nil {
		return nil, errwrap.Wrapf("could not start transaction: {{err}}", err)
	}

	return txn, nil
}

// deferredRollback can be used to rollback a transaction in a defer.
// It will log an error if it fails
func deferredRollback(txn *sql.Tx) {
	err := txn.Rollback()
	switch {
	case err == sql.ErrTxDone:
		// transaction has already been committed or rolled back
		log.Printf("[DEBUG]: %v", err)
	case err != nil:
		log.Printf("[ERR] could not rollback transaction: %v", err)
	}
}
