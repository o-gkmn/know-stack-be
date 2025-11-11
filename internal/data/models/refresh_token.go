package models

import "time"

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"index"`
	ExpiresAt time.Time `gorm:""`
	IsRevoked bool      `gorm:"default:false"`
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
