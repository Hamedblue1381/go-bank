package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/HamedBlue1381/hamed-bank/api"
	db "github.com/HamedBlue1381/hamed-bank/db/bankmodel"
	"github.com/HamedBlue1381/hamed-bank/doc"
	"github.com/HamedBlue1381/hamed-bank/gapi"
	"github.com/HamedBlue1381/hamed-bank/pb"
	"github.com/HamedBlue1381/hamed-bank/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
	go runGatewayServer(config, store)
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

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create the server", err)
	}
	jsonOptions :=
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		})

	grpcMux := runtime.NewServeMux(jsonOptions)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterHamedBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	swaggerHandler := http.FileServer(http.FS(doc.SwaggerFiles))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create the listener", err)
	}

	log.Printf("Starting the HTTP Gateway server: %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start the HTTP gateway server", err)
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
