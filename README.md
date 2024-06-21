# go-ratelimiter

[![build](https://github.com/adobromilskiy/go-ratelimiter/actions/workflows/test.yml/badge.svg)](https://github.com/adobromilskiy/go-ratelimiter/actions/workflows/test.yml)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/adobromilskiy/go-ratelimiter)](https://pkg.go.dev/github.com/adobromilskiy/go-ratelimiter)
[![Go Report Card](https://goreportcard.com/badge/github.com/adobromilskiy/go-ratelimiter)](https://goreportcard.com/report/github.com/adobromilskiy/go-ratelimiter)
[![Coverage Status](https://coveralls.io/repos/github/adobromilskiy/go-ratelimiter/badge.svg?branch=main&kill_cache=1)](https://coveralls.io/github/adobromilskiy/go-ratelimiter?branch=main)

Simple ratelimiter on Golang.

RateLimiter works with Redis.

To install go-ratelimiter, use `go get`:

```bash
go get -u github.com/adobromilskiy/go-ratelimiter
```

## Quickstart

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/adobromilskiy/go-ratelimiter"
	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	rl := ratelimiter.NewRateLimiter(rdb)

	userID := "user123"
	key := fmt.Sprintf("{user:%s}:msg", userID)
	duration := time.Second
	limit := int64(5)

	for i := 0; i < 10; i++ {
		ok, err := rl.Do(ctx, key, limit, duration)
		if err != nil {
			fmt.Println(err)
			return
		}
		if !ok {
			fmt.Println("limit exceeded")
			return
		}
		fmt.Println("message sent")
	}
}
```

Also you can use Redis Cluster:

```go
addrs := []string{"172.18.0.2:6379", "172.18.0.3:6379", "172.18.0.4:6379", "172.18.0.5:6379", "172.18.0.6:6379", "172.18.0.7:6379"}

rdb := redis.NewClusterClient(&redis.ClusterOptions{
    Addrs: addrs,
})

rl := ratelimiter.NewRateLimiter(rdb)
```

## Package overview

```
package ratelimiter // import "github.com/adobromilskiy/go-ratelimiter"


TYPES

type RateLimiter struct {
        client RedisRateLimiter
}

func NewRateLimiter(client RedisRateLimiter) *RateLimiter

func (r *RateLimiter) Do(ctx context.Context, key string, limit int64, duration time.Duration) (bool, error)

type RedisRateLimiter interface {
        Pipeline() redis.Pipeliner
}
```