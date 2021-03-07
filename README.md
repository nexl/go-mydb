## mydb

Mydb is a small and lightweight library written in Go for automatically route read-only queries to read replicas, and all other queries to the master DB ( master - slave replication)

## Run test
```
go test -race ./...
```