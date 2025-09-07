package redisx

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Producer struct {
	rdb *redis.Client
}

func NewProducer(rdb *redis.Client) *Producer { return &Producer{rdb: rdb} }

func (p *Producer) Publish(ctx context.Context, stream string, fields map[string]interface{}) (string, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	return p.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: fields,
	}).Result()
}
