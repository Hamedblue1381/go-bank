package main

import (
	"database/sql"
	"log"

	"github.com/HamedBlue1381/go-postgres-gRPC/api"
	db "github.com/HamedBlue1381/go-postgres-gRPC/db/bankmodel"
	"github.com/HamedBlue1381/go-postgres-gRPC/util"
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
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create the server", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start the server", err)
	}
}
