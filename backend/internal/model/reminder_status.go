package model

const (
	ReminderPending    = "pending"
	ReminderProcessing = "processing"
	ReminderRetrying   = "retrying"
	ReminderSent       = "sent"
	ReminderFailed     = "failed"
)

// Failure type classification for reminder execution failures.
const (
	FailureWahaTimeout    = "waha_timeout"
	FailureInvalidTarget  = "invalid_target"
	FailureDBLockFail     = "db_lock_fail"
	FailureUnknownSendErr = "unknown_send_error"
)
