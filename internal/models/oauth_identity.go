package models

import "time"

type OAuthIdentity struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"index;not null"`
	Provider    string    `gorm:"index;size:50;not null"`   // "google"
	ProviderUID string    `gorm:"index;size:255;not null"`  // sub
	Email       string    `gorm:"size:255"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (OAuthIdentity) TableName() string { return "oauth_identities" }
