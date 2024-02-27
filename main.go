package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"sgithub.com/techschool/simplebank/api"
	db "sgithub.com/techschool/simplebank/db/sqlc"
	"sgithub.com/techschool/simplebank/util"
)

func main() {
	config, err := util.LoadConfigDB_Server(".")
	if err != nil {
		log.Fatal("can not load the config file", err)
	}
	conn, err := sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatal("can not connect to the DB", err)
	}
	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("can not initiate", err)
	}
	err = server.Start(config.ServerAddress) // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Fatal("can not start server", err)
	}

}
