package queue

import (
	"nemoris/internal/redis"
)

const execQueueKey = "reminder:queue:exec"

// Enqueue adds reminderID as an execution marker to the queue.
// Returns redis.ErrRedisUnavailable when Redis is not connected.
func Enqueue(reminderID string) error {
	return redis.ListPush(execQueueKey, reminderID)
}
