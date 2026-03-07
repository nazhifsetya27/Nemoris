package model

import "time"

type Message struct {
	BaseModel
	From      string    `gorm:"index:idx_message_duplicate"`
	Body      string    `gorm:"index:idx_message_duplicate"`
	CreatedAt time.Time `gorm:"index:idx_message_duplicate"`
}