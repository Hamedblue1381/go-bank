package gapi

import (
	"fmt"

	db "github.com/HamedBlue1381/hamed-bank/db/bankmodel"
	"github.com/HamedBlue1381/hamed-bank/pb"
	"github.com/HamedBlue1381/hamed-bank/token"
	"github.com/HamedBlue1381/hamed-bank/util"
	"github.com/HamedBlue1381/hamed-bank/worker"
)

type Server struct {
	pb.UnimplementedHamedBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cant create token maker: %w", err)
	}
	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
