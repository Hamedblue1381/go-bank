package worker

import (
	"context"

	db "github.com/HamedBlue1381/hamed-bank/db/bankmodel"
	"github.com/hibiken/asynq"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RediTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRediTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{},
	)

	return &RediTaskProcessor{
		server: server,
		store:  store,
	}
}
func (processor *RediTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}
