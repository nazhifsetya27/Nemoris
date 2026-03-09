package queue

import (
	"fmt"

	"nemoris/internal/redis"
	"nemoris/internal/utils"
)

// TryLock attempts to acquire execution lock for reminderID.
// Returns (true, nil) when acquired, (false, nil) when already held.
// Returns (false, err) when Redis unavailable or operation fails.
func TryLock(reminderID string) (acquired bool, err error) {
	acquired, err = redis.TryLockStrict(reminderID)
	if err != nil {
		utils.LogSchedulerWarn(fmt.Sprintf("queue_fail type=lock_timeout reminder_id=%s err=%v", reminderID, err))
		return false, err
	}
	return acquired, nil
}

// Release releases the execution lock for reminderID.
// Returns redis.ErrRedisUnavailable when Redis is not connected.
func Release(reminderID string) error {
	err := redis.ReleaseStrict(reminderID)
	if err != nil {
		utils.LogSchedulerWarn(fmt.Sprintf("queue_fail type=lock_timeout reminder_id=%s err=%v", reminderID, err))
		return err
	}
	return nil
}
