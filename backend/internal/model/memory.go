package model

import "time"

type Memory struct {
	BaseModel
	From      string    `gorm:"column:from;type:text"`
	Content   string    `gorm:"type:text"`
	Lang      string    `gorm:"type:text"`
	CreatedAt time.Time
}
