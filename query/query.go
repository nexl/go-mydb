package query

import (
	"context"
	"database/sql"
	"log"

	"../connection"
	"../mydb"
)

// Query returns rows and error from query string
func Query(db *mydb.DB, query string, args ...interface{}) (*sql.Rows, error) {
	var rows *sql.Rows
	var err error

	ch := make(chan struct{})
	go func() {
		rows, err = connection.ReadReplicaRoundRobin(db).Query(query, args...)
		if err != nil {
			log.Printf("Something went wrong, '%s'", err)
		}
		close(ch)
	}()
	<-ch
	return rows, err
}

func QueryContext(db *mydb.DB, ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return connection.ReadReplicaRoundRobin(db).QueryContext(ctx, query, args...)
}

func QueryRow(db *mydb.DB, query string, args ...interface{}) *sql.Row {
	return connection.ReadReplicaRoundRobin(db).QueryRow(query, args...)
}

func QueryRowContext(db *mydb.DB, ctx context.Context, query string, args ...interface{}) *sql.Row {
	return connection.ReadReplicaRoundRobin(db).QueryRowContext(ctx, query, args...)
}

func Begin(db *mydb.DB) (*sql.Tx, error) {
	return db.Master.Begin()
}

func BeginTx(db *mydb.DB, ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return db.Master.BeginTx(ctx, opts)
}

func Exec(db *mydb.DB, query string, args ...interface{}) (sql.Result, error) {
	return db.Master.Exec(query, args...)
}

func ExecContext(db *mydb.DB, ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.Master.ExecContext(ctx, query, args...)
}

// Prepare is used for the preparation of the statement
func Prepare(db *mydb.DB, query string) (*sql.Stmt, error) {
	var stmt *sql.Stmt
	var err error

	ch := make(chan struct{})
	go func() {
		stmt, err = db.Master.Prepare(query)
		if err != nil {
			log.Printf("Something went wrong, '%s'", err)
		}
		close(ch)
	}()
	<-ch
	return stmt, err
}

func PrepareContext(db *mydb.DB, ctx context.Context, query string) (*sql.Stmt, error) {
	return db.Master.PrepareContext(ctx, query)
}
