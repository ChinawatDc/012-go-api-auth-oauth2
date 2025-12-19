package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"uniqueIndex;size:255;not null"`
	Name      string    `gorm:"size:255"`
	AvatarURL string    `gorm:"size:500"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
