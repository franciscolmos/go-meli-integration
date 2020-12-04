package model

import (
	"time"
)

type Item struct {
	ID        uint           `gorm:"primaryKey"`
	Title string
	Quantity int
	Price float64
	FirstPicture string
	ItemId string
	CreatedAt time.Time
	UpdatedAt time.Time
}
