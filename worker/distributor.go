package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error
}
type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpts asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpts)
	return &RedisTaskDistributor{
		client: client,
	}
}
