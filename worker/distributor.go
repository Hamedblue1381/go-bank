package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDirtributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

type RedisTaskDirtributor struct {
	client *asynq.Client
}

func newRedisTaskDirtributor(redisOpt asynq.RedisClientOpt) TaskDirtributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDirtributor{
		client: client,
	}
}
