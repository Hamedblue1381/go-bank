package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/HamedBlue1381/hamed-bank/api"
	db "github.com/HamedBlue1381/hamed-bank/db/bankmodel"
	"github.com/HamedBlue1381/hamed-bank/gapi"
	"github.com/HamedBlue1381/hamed-bank/pb"
	"github.com/HamedBlue1381/hamed-bank/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create the server", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterHamedBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create the listener", err)
	}

	log.Printf("Starting the gRPC server: %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start the gRPC server", err)
	}
}
func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create the server", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start the server", err)
	}
}
