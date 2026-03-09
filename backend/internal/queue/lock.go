package queue

import (
	"nemoris/internal/redis"
)

// TryLock attempts to acquire execution lock for reminderID.
// Returns (true, nil) when acquired, (false, nil) when already held.
// Returns (false, err) when Redis unavailable or operation fails.
func TryLock(reminderID string) (acquired bool, err error) {
	return redis.TryLockStrict(reminderID)
}

// Release releases the execution lock for reminderID.
// Returns redis.ErrRedisUnavailable when Redis is not connected.
func Release(reminderID string) error {
	return redis.ReleaseStrict(reminderID)
}
