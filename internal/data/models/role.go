package models

import "time"

type Role struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"unique not null"`
	IsDefault bool      `gorm:"default:false"`
	Claims    []Claim   `gorm:"many2many:role_claims;"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (Role) TableName() string {
	return "roles"
}
