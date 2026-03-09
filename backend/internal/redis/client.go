package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"nemoris/internal/config"
	"nemoris/internal/utils"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

// ErrRedisUnavailable is returned when Redis client is nil (not connected).
var ErrRedisUnavailable = errors.New("redis unavailable")

const (
	execLockPrefix = "reminder:exec:"
	execLockTTL    = 30 * time.Second
)

// Init connects to Redis when REDIS_HOST is set. If unset, Redis stays disabled.
// Redis failure at init does not exit; client remains nil and callers fallback.
func Init() {
	if config.App.RedisHost == "" {
		utils.LogScheduler("Redis disabled (REDIS_HOST empty)")
		return
	}

	addr := fmt.Sprintf("%s:%s", config.App.RedisHost, config.App.RedisPort)
	client = redis.NewClient(&redis.Options{
		Addr:         addr,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		utils.LogSchedulerWarn("Redis connect failed, reminder execution will not use lock: " + err.Error())
		client = nil
		return
	}

	utils.LogScheduler("Redis connected for execution lock")
}

// TryLock attempts to acquire an execution lock for reminderID.
// Returns (true, nil) when lock acquired or when Redis is unavailable (fallback).
// Returns (false, nil) when Redis is up but lock already held (another worker).
func TryLock(reminderID string) (acquired bool, err error) {
	if client == nil {
		return true, nil
	}

	key := execLockPrefix + reminderID
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ok, err := client.SetNX(ctx, key, "1", execLockTTL).Result()
	if err != nil {
		utils.LogSchedulerWarn("Redis lock failed, proceeding without lock: " + err.Error())
		return true, nil
	}
	return ok, nil
}

// Release releases the execution lock for reminderID. No-op when Redis unavailable.
func Release(reminderID string) {
	if client == nil {
		return
	}

	key := execLockPrefix + reminderID
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_ = client.Del(ctx, key).Err()
}

// TryLockStrict attempts to acquire execution lock. Returns ErrRedisUnavailable when client nil.
// Returns (false, nil) when lock already held. Used by queue package; no fallback.
func TryLockStrict(reminderID string) (acquired bool, err error) {
	if client == nil {
		return false, ErrRedisUnavailable
	}
	key := execLockPrefix + reminderID
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	ok, err := client.SetNX(ctx, key, "1", execLockTTL).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

// ReleaseStrict releases execution lock. Returns ErrRedisUnavailable when client nil.
func ReleaseStrict(reminderID string) error {
	if client == nil {
		return ErrRedisUnavailable
	}
	key := execLockPrefix + reminderID
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return client.Del(ctx, key).Err()
}

// ListPush pushes value onto the right of list key. Returns ErrRedisUnavailable when client nil.
func ListPush(key, value string) error {
	if client == nil {
		return ErrRedisUnavailable
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return client.RPush(ctx, key, value).Err()
}

// ListPop pops value from the left of list key. Returns ErrRedisUnavailable when client nil.
// Returns redis.Nil when list empty.
func ListPop(key string) (string, error) {
	if client == nil {
		return "", ErrRedisUnavailable
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return client.LPop(ctx, key).Result()
}
