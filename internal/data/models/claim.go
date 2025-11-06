package models

import "time"

type Claim struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"unique not null"`
	Roles     []Role    `gorm:"many2many:role_claims;"`
	Users     []User    `gorm:"many2many:user_claims;"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (Claim) TableName() string {
	return "claims"
}