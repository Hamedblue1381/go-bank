package main

import (
	"database/sql"
	"log"

	"github.com/HamedBlue1381/go-postgres-gRPC/api"
	db "github.com/HamedBlue1381/go-postgres-gRPC/db/bankmodel"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/hamed_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8000"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start the server", err)
	}
}
