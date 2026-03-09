package redis

import (
	"context"
	"fmt"
	"time"

	"nemoris/internal/config"
	"nemoris/internal/utils"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

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
