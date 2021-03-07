package healthcheck

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

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
		log.Printf("Something went wrong when creating healthcheck mock database, %s\n", err)
	}
}

func TestPingMaster(t *testing.T) {
	db := mydb.NewDB(master, readReplicaOne, readReplicaTwo, readReplicaThree)
	err := PingMaster(db)
	assert.Nil(t, err)
}

func TestPingReplica(t *testing.T) {
	db := mydb.NewDB(master, readReplicaOne, readReplicaTwo, readReplicaThree)
	for i := range db.Readreplicas {
		err := PingReplica(db.Readreplicas[i])
		assert.Nil(t, err)
	}
}

func TestPingContext(t *testing.T) {
	db := mydb.NewDB(master, readReplicaOne, readReplicaTwo, readReplicaThree)
	backgroundContext := context.Background()
	context2SecondsTimeout, cancelFunc2SecondsTimeout := context.WithTimeout(backgroundContext, time.Second*2)
	defer cancelFunc2SecondsTimeout()
	PingContext(context2SecondsTimeout, db)
}
