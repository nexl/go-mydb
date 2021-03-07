package connection

import (
	"database/sql"
	"log"
	"testing"

	"../mydb"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var master, readReplicaOne, readReplicaTwo, readReplicaThree *sql.DB

func init() {
	var err error
	master, _, err = sqlmock.New()
	readReplicaOne, _, err = sqlmock.New()
	readReplicaTwo, _, err = sqlmock.New()
	readReplicaThree, _, err = sqlmock.New()
	if err != nil {
		log.Printf("Something went wrong when creating connection mock database, %s\n", err)
	}
}

func TestReadReplicaRoundRobin(t *testing.T) {
	db := mydb.NewDB(master, readReplicaOne, readReplicaTwo, readReplicaThree)
	err := ReadReplicaRoundRobin(db)
	assert.NotNil(t, err)
	assert.Equal(t, err, db.Readreplicas[1])
}

func TestSetMaxIdleConns(t *testing.T) {
	db := mydb.NewDB(master, readReplicaOne, readReplicaTwo, readReplicaThree)
	SetMaxIdleConns(db, 1)
	stats := db.Master.Stats()
	assert.Equal(t, stats.Idle, 1, "Max idle connection is 1")
}

func TestMaxOpenConnections(t *testing.T) {
	mydatabase := mydb.NewDB(master, readReplicaOne, readReplicaTwo, readReplicaThree)
	SetMaxOpenConns(mydatabase, 5)
	stats := mydatabase.Master.Stats()
	assert.Equal(t, stats.MaxOpenConnections, 5, "Max connection is 5")
}
