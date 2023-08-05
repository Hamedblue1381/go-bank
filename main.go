package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/HamedBlue1381/hamed-bank/api"
	db "github.com/HamedBlue1381/hamed-bank/db/bankmodel"
	"github.com/HamedBlue1381/hamed-bank/doc"
	"github.com/HamedBlue1381/hamed-bank/gapi"
	"github.com/HamedBlue1381/hamed-bank/mail"
	"github.com/HamedBlue1381/hamed-bank/pb"
	"github.com/HamedBlue1381/hamed-bank/util"
	"github.com/HamedBlue1381/hamed-bank/worker"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config:")
	}
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Msg("cannot connect to db:")
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{

		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProcessor(config, redisOpt, store)
	go runGatewayServer(config, store, taskDistributor)
	runGrpcServer(config, store, taskDistributor)
}
func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Msg("cannot create migration instance")
	}
	if err = migration.Up(); err != nil {
		log.Fatal().Msg("failed to run the migrate up:")
	}

	log.Info().Msg("db migrated successfully")
}

func runTaskProcessor(config util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRediTaskProcessor(redisOpt, store, mailer)
	log.Info().Msg("starting task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}

func runGrpcServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg("cannot create the server")
	}
	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterHamedBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create the listener")
	}

	log.Info().Msgf("Starting the gRPC server: %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Msg("cannot start the gRPC server")
	}
}

func runGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg("cannot create the server")
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
		log.Fatal().Msg("cannot register handler server:")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	swaggerHandler := http.FileServer(http.FS(doc.SwaggerFiles))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create the listener")
	}

	log.Info().Msgf("Starting the HTTP Gateway server: %s", listener.Addr().String())

	handler := gapi.HttpLogger(mux)

	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Msg("cannot start the HTTP gateway server")
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("cannot create the server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot start the server")
	}
}
