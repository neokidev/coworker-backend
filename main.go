package main

import (
	"database/sql"
	"github.com/ot07/coworker-backend/api"
	db "github.com/ot07/coworker-backend/db/sqlc"
	"github.com/ot07/coworker-backend/util"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start to server:", err)
	}
}
