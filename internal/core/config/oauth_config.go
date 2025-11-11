package config

import (
	"knowstack/internal/utils"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuth struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

func DefaultOAuthConfigFromEnv() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     utils.GetEnv("GOOGLE_CLIENT_ID", ""),
		ClientSecret: utils.GetEnv("GOOGLE_CLIENT_SECRET", ""),
		RedirectURL:  utils.GetEnv("GOOGLE_REDIRECT_URL", ""),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
