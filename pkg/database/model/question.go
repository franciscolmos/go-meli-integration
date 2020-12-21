package model

import (
	"time"
)

type Question struct {
	ID        uint           `gorm:"primaryKey"`
	Text string
	ItemTitle string
	CreatedAt time.Time
	UpdatedAt time.Time
}
