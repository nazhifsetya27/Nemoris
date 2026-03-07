package model

import "time"

type Reminder struct {
	BaseModel
	From      string
	Task      string
	RawTime   string
	RemindAt  time.Time `gorm:"index"`
	CreatedAt time.Time
}