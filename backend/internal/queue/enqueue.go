package queue

import (
	"time"

	"nemoris/internal/redis"
)

const (
	execQueueKey    = "reminder:queue:exec"
	delayedQueueKey = "reminder:queue:delayed"
)

// Enqueue adds reminderID as an execution marker to the queue.
// Returns redis.ErrRedisUnavailable when Redis is not connected.
func Enqueue(reminderID string) error {
	return redis.ListPush(execQueueKey, reminderID)
}

// EnqueueDelayed adds reminderID as delayed execution marker. executeAt is when it becomes eligible.
// Returns redis.ErrRedisUnavailable when Redis is not connected.
func EnqueueDelayed(reminderID string, executeAt time.Time) error {
	return redis.SortedSetAdd(delayedQueueKey, float64(executeAt.Unix()), reminderID)
}

// PromoteDelayedToExec moves due items from delayed set to exec list. Best-effort; logs on failure.
func PromoteDelayedToExec() {
	now := float64(time.Now().Unix())
	members, err := redis.SortedSetRangeByScore(delayedQueueKey, now, 50)
	if err != nil {
		return
	}
	for _, id := range members {
		_ = redis.ListPush(execQueueKey, id)
		_ = redis.SortedSetRem(delayedQueueKey, id)
	}
}
