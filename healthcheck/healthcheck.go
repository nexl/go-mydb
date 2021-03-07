package healthcheck

import (
	"context"
	"database/sql"
	"log"

	"../mydb"
)

// PingMaster returns nil if master database is alive
// returns an error otherwise
func PingMaster(db *mydb.DB) error {
	if err := db.Master.Ping(); err != nil {
		log.Println("Master db is not available")
		return err
	}
	return nil
}

// PingReplica returns nil if replica database is alive
// returns an error otherwise
func PingReplica(db interface{}) error {
	if err := db.(*sql.DB).Ping(); err != nil {
		log.Println(db, " is not available")
		return err
	}
	return nil
}

// PingContext returns nil if both master and replicas are alive
// returns an error otherwise
func PingContext(ctx context.Context, db *mydb.DB) error {
	if err := db.Master.PingContext(ctx); err != nil {
		log.Println("Master context db is not available")
		return err
	}
	for i := range db.Readreplicas {
		if err := db.Readreplicas[i].(*sql.DB).PingContext(ctx); err != nil {
			log.Println(db, "context is not available")
			return err
		}
	}
	return nil
}
