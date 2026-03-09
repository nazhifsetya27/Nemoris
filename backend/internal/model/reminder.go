package model

import "time"

type Reminder struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	From        string
	Task        string
	RawTime     string
	RemindAt    time.Time
	Status      string
	SentAt      *time.Time
	RetryCount  int
	LastError   string `gorm:"type:text"`
	FailureType string `gorm:"column:failure_type;size:32"`
	CreatedAt   time.Time
}