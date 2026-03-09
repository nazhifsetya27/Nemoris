package queue

import (
	"fmt"
	"time"

	"nemoris/internal/redis"
	"nemoris/internal/utils"
)

const (
	execQueueKey    = "reminder:queue:exec"
	delayedQueueKey = "reminder:queue:delayed"
)

// Enqueue adds reminderID as an execution marker to the queue.
// Returns redis.ErrRedisUnavailable when Redis is not connected.
func Enqueue(reminderID string) error {
	err := redis.ListPush(execQueueKey, reminderID)
	if err != nil {
		utils.LogSchedulerWarn(fmt.Sprintf("queue_fail type=enqueue_fail err=%v", err))
		return err
	}
	return nil
}

// EnqueueDelayed adds reminderID as delayed execution marker. executeAt is when it becomes eligible.
// Returns redis.ErrRedisUnavailable when Redis is not connected.
func EnqueueDelayed(reminderID string, executeAt time.Time) error {
	err := redis.SortedSetAdd(delayedQueueKey, float64(executeAt.Unix()), reminderID)
	if err != nil {
		utils.LogSchedulerWarn(fmt.Sprintf("queue_fail type=enqueue_fail err=%v", err))
		return err
	}
	return nil
}

// PromoteDelayedToExec moves due items from delayed set to exec list. Best-effort; logs on failure.
func PromoteDelayedToExec() {
	now := float64(time.Now().Unix())
	members, err := redis.SortedSetRangeByScore(delayedQueueKey, now, 50)
	if err != nil {
		utils.LogSchedulerWarn(fmt.Sprintf("queue_fail type=dequeue_fail err=%v", err))
		return
	}
	for _, id := range members {
		_ = redis.ListPush(execQueueKey, id)
		_ = redis.SortedSetRem(delayedQueueKey, id)
	}
}
