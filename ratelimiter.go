package ratelimiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRateLimiter interface {
	Pipeline() redis.Pipeliner
}

type RateLimiter struct {
	client RedisRateLimiter
}

func NewRateLimiter(client RedisRateLimiter) *RateLimiter {
	return &RateLimiter{client: client}
}

func (r *RateLimiter) Do(ctx context.Context, key string, limit int64, duration time.Duration) (bool, error) {
	pipe := r.client.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, duration)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	if incr.Val() > limit {
		return false, nil
	}

	return true, nil
}
