package model

import "time"

type Message struct {
	ID        uint      `gorm:"primaryKey"`
	From      string
	Body      string
	CreatedAt time.Time
}