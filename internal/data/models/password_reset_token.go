package models

import "time"

type PasswordResetToken struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	IsUsed    bool      `gorm:"default:false"`
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}
