package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"sgithub.com/techschool/simplebank/util"
)

var testQueries *Queries
var testDB *sql.DB
var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfigDB_Server("../../")
	if err != nil {
		log.Fatal("can not load the config file", err)
	}
	testDB, err = sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatal("cannot connect to the DB because of", err)
		return
	}
	testQueries = New(testDB)
	testStore = NewStore(testDB)
	os.Exit(m.Run())
}
