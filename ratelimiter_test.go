package ratelimiter

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
)

func TestDo(t *testing.T) {
	db, mock := redismock.NewClientMock()

	ctx := context.Background()
	ip := "127.0.0.1"
	limit := int64(10)
	duration := time.Minute

	key := fmt.Sprintf("{ip:%s}:requests", ip)

	rl := NewRateLimiter(db)

	mock.ExpectIncr(key).SetVal(1)
	mock.ExpectExpire(key, duration).SetVal(true)

	ok, err := rl.Do(ctx, key, limit, duration)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ok {
		t.Fatalf("unexpected rate limit")
	}
}

func TestDo_RateLimit(t *testing.T) {
	db, mock := redismock.NewClientMock()

	ctx := context.Background()
	ip := "127.0.0.1"
	limit := int64(10)
	duration := time.Minute

	key := fmt.Sprintf("{ip:%s}:requests", ip)

	rl := NewRateLimiter(db)

	mock.ExpectIncr(key).SetVal(11)
	mock.ExpectExpire(key, duration).SetVal(false)

	ok, err := rl.Do(ctx, key, limit, duration)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ok {
		t.Fatalf("expected rate limit")
	}
}

func TestDo_Error(t *testing.T) {
	db, mock := redismock.NewClientMock()

	ctx := context.Background()
	ip := "127.0.0.1"
	limit := int64(10)
	duration := time.Minute

	key := fmt.Sprintf("{ip:%s}:requests", ip)

	rl := NewRateLimiter(db)

	mock.ExpectIncr(key).SetVal(11)
	mock.ExpectExpire(key, duration).SetErr(fmt.Errorf("error"))

	ok, err := rl.Do(ctx, key, limit, duration)
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}

	if ok {
		t.Fatalf("expected rate limit")
	}
}
