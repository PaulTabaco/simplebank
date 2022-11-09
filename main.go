package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"parus.i234.me/paultabaco/simplebank/api"
	db "parus.i234.me/paultabaco/simplebank/db/sqlc"
	"parus.i234.me/paultabaco/simplebank/util"
)

func main() {
	config, err := util.LoadConfig(".") // "." - current folder
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can not connect to db", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
