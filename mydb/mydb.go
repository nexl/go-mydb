package mydb

import (
	"database/sql"
)

// DB represent database struct
type DB struct {
	Master       *sql.DB
	Readreplicas []interface{}
	Count        int
}

// NewDB returns a new DB (master and replicas)
func NewDB(master *sql.DB, readreplicas ...interface{}) *DB {
	return &DB{
		Master:       master,
		Readreplicas: readreplicas,
	}
}
