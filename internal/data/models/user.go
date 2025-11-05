package models

import (
	"knowstack/internal/utils"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"unique"`
	Email     string    `gorm:"unique"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Password = utils.HashPassword(u.Password)
	return
}


func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.Password = utils.HashPassword(u.Password)
	return
}