package model

import "time"

type Reminder struct {
	ID        uint      `gorm:"primaryKey"`
	From      string
	Task      string
	RawTime   string
	RemindAt  time.Time
	CreatedAt time.Time
}