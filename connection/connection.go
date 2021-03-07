package connection

import (
	"database/sql"
	"log"
	"time"

	"../healthcheck"
	"../mydb"
)

// ReadReplicaRoundRobin returns available replica database
// maintenance of read replicas is always performed on one replica at a time,
// there's no chance if all read replicas are not available
func ReadReplicaRoundRobin(db *mydb.DB) *sql.DB {
	var availableReadReplicas []interface{}
	ch := make(chan struct{})
	go func() {
		db.Count++
		for i := range db.Readreplicas {
			currentReplica := db.Readreplicas[i]
			if err := healthcheck.PingReplica(currentReplica); err == nil {
				availableReadReplicas = append(availableReadReplicas, currentReplica)
			}
		}
		close(ch)
	}()
	<-ch
	index := db.Count % len(availableReadReplicas)
	// Round robin distribution always starts from index 1 -> 2 -> 0
	log.Println("Using read replica #", index)
	return availableReadReplicas[index].(*sql.DB)
}

// Close all database (master and read replicas)
func Close(db *mydb.DB) error {
	db.Master.Close()
	for i := range db.Readreplicas {
		if err := db.Readreplicas[i].(*sql.DB).Close(); err != nil {
			log.Printf("Failed to close read replica #%d db %s \n", i, err)
			return err
		}
	}
	return nil
}

// SetConnMaxLifetime set max length of time that a connection can be reused for
func SetConnMaxLifetime(db *mydb.DB, d time.Duration) {
	db.Master.SetConnMaxLifetime(d)
	for i := range db.Readreplicas {
		if err := healthcheck.PingReplica(db.Readreplicas[i]); err != nil {
			db.Readreplicas[i].(*sql.DB).SetConnMaxLifetime(d)
		} else {
			log.Printf("Failed to SetConnMaxLifetime read replica # %d db, %s \n", i, err)
		}
	}
}

// SetMaxIdleConns set max idle connections to be retained in the connection pool
// (default value is 2)
func SetMaxIdleConns(db *mydb.DB, n int) {
	db.Master.SetMaxIdleConns(n)
	for i := range db.Readreplicas {
		if err := healthcheck.PingReplica(db.Readreplicas[i]); err != nil {
			db.Readreplicas[i].(*sql.DB).SetMaxIdleConns(n)
		} else {
			log.Printf("Failed to SetMaxIdleConns read replica # %d db, %s \n", i, err)
		}
	}
}

// SetMaxOpenConns set limit on the number of open connection at the same time
func SetMaxOpenConns(db *mydb.DB, n int) {
	db.Master.SetMaxOpenConns(n)
	for i := range db.Readreplicas {
		if err := healthcheck.PingReplica(db.Readreplicas[i]); err != nil {
			db.Readreplicas[i].(*sql.DB).SetMaxOpenConns(n)
		} else {
			log.Printf("Failed to SetMaxOpenConns read replica # %d db, %s \n", i, err)
		}
	}
}
