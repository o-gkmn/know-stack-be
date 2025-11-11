package models

import (
	"knowstack/internal/utils"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Username     string    `gorm:"unique"`
	Email        string    `gorm:"unique"`
	Password     string    `gorm:""`
	RoleID       uint      `gorm:"not null"`
	Role         Role      `gorm:"foreignKey:RoleID"`
	Claims       []Claim   `gorm:"many2many:user_claims;"`
	GoogleID     string    `gorm:"uniqueIndex"`
	Provider     string    `gorm:"default:local"`
	ProfileImage string    `gorm:""`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Provider == "local" {
		u.Password = utils.HashPassword(u.Password)
	}
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if u.Provider == "" && tx.Statement.Changed("Password") {
		u.Password = utils.HashPassword(u.Password)
	}
	return
}
