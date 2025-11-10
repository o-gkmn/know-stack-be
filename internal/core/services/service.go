package services

import "gorm.io/gorm"

type Service struct {
	UserService *UserService
	ClaimService *ClaimService
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		UserService: NewUserService(db),
		ClaimService: NewClaimService(db),
	}
}
