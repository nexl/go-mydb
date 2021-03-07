package query

import (
	"database/sql"
	"log"
	"regexp"
	"testing"

	"../mydb"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var master, readReplicaOne, readReplicaTwo, readReplicaThree *sql.DB
var mockMaster, mockOne, mockTwo, mockThree sqlmock.Sqlmock

func init() {
	var err error
	master, mockMaster, err = sqlmock.New()
	readReplicaOne, mockOne, err = sqlmock.New()
	readReplicaTwo, mockTwo, err = sqlmock.New()
	readReplicaThree, mockThree, err = sqlmock.New()
	if err != nil {
		log.Printf("Something went wrong when creating query mock database, %s\n", err)
	}
}

func TestQuery(t *testing.T) {
	db := mydb.NewDB(master, readReplicaOne, readReplicaTwo, readReplicaThree)
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "one")

	mockTwo.ExpectQuery(regexp.QuoteMeta("select * from user")).WithArgs("1").WillReturnRows(rows)

	rs, err := Query(db, "select * from user", "1")
	columns, _ := rs.Columns()

	assert.Nil(t, err)
	assert.Equal(t, len(columns), 2, "Total column length is 2")
	assert.Equal(t, columns[0], "id", "First column name is id")
	assert.Equal(t, columns[1], "name", "Second column name is 'name'")
}

func TestPrepare(t *testing.T) {
	sampleQuery := "INSERT into user values "
	db := mydb.NewDB(master, readReplicaOne, readReplicaTwo, readReplicaThree)

	mockMaster.ExpectPrepare(sampleQuery)
	Prepare(db, sampleQuery)

	err := mockMaster.ExpectationsWereMet()
	assert.Nil(t, err)
}
