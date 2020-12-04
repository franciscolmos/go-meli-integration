package model

import (
	"time"
)

type Question struct {
	ID        uint           `gorm:"primaryKey"`
	Text string
	Question_Id int
	ItemTitle string
	CreatedAt time.Time
	UpdatedAt time.Time
}
