package model

import (
	"time"
)

type User struct {
	ID        uint           `gorm:"primaryKey"`
	AccessToken string
	RefreshToken string
	UserIdMeli int
	CreatedAt time.Time
	UpdatedAt time.Time
}
