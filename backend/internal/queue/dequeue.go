package queue

import (
	"nemoris/internal/redis"
)

// Dequeue removes and returns the next execution marker from the queue.
// Returns ("", err) when queue is empty or Redis unavailable.
func Dequeue() (reminderID string, err error) {
	return redis.ListPop(execQueueKey)
}
