package queue

import (
	"errors"
	"fmt"

	"nemoris/internal/redis"
	"nemoris/internal/utils"
)

// Dequeue removes and returns the next execution marker from the queue.
// Returns ("", err) when queue is empty or Redis unavailable.
func Dequeue() (reminderID string, err error) {
	reminderID, err = redis.ListPop(execQueueKey)
	if err != nil && !errors.Is(err, redis.Nil) {
		utils.LogSchedulerWarn(fmt.Sprintf("queue_fail type=dequeue_fail err=%v", err))
		return "", err
	}
	return reminderID, err
}
