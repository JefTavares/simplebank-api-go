package main

import (
	"database/sql"
	"log"

	"github.com/jeftavares/simple_bank/api"
	db "github.com/jeftavares/simple_bank/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	dbDriver     = "postgres"
	dbSource     = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	severAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(severAddress)

	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
