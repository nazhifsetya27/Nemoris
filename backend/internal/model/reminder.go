package model

import "time"

type Reminder struct {
	ID         string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	From       string
	Task       string
	RawTime    string
	RemindAt   time.Time
	Status     string
	SentAt     *time.Time
	RetryCount int
	CreatedAt  time.Time
}