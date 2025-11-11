package services

import (
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type Service struct {
	UserService  *UserService
	ClaimService *ClaimService
	OAuthService *OAuthService
}

func NewService(db *gorm.DB, oauthConfig *oauth2.Config) *Service {
	return &Service{
		UserService:  NewUserService(db),
		ClaimService: NewClaimService(db),
		OAuthService: NewOAuthService(db, oauthConfig),
	}
}
