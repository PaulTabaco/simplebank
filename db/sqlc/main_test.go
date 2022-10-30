package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// / TODO:- get these vars from env
const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:mysecretpassword@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("can not connect to db", err)
	}
	testQueries = New(conn)

	os.Exit(m.Run())
}
