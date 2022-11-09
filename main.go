package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"parus.i234.me/paultabaco/simplebank/api"
	db "parus.i234.me/paultabaco/simplebank/db/sqlc"
)

const (
	dbDriver     = "postgres"
	dbSource     = "postgresql://root:mysecretpassword@localhost:5432/simple_bank?sslmode=disable"
	serverAdress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("can not connect to db", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAdress)

	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
