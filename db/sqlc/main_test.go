package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/jeftavares/simple_bank/util"
	_ "github.com/lib/pq"
)

//sera utiizado o arquivo app.env
// const (
// 	dbDriver = "postgres"
// 	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
// )

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	//criar o testDB uma unica vez, passar = e não atribuir :=
	config, err := util.LoadConfig("../..") //volta 2 pastas rss
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	//utilizado o aqruivo app.env testDB, err = sql.Open(dbDriver, dbSource)
	//Exemplo utilizando pgx
	//ctx := context.Background()
	//conn, err := pgx.Connect(ctx, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
