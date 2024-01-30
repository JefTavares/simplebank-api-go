package main

import (
	"database/sql"
	"log"

	"github.com/jeftavares/simple_bank/api"
	db "github.com/jeftavares/simple_bank/db/sqlc"
	"github.com/jeftavares/simple_bank/util"

	_ "github.com/lib/pq"
)

//Utilizado o viper com o arquivo app.env
// const (
// 	dbDriver     = "postgres"
// 	dbSource     = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
// 	severAddress = "0.0.0.0:8080"
// )

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
		log.Fatal("cannot start server:", err)
	}
}
