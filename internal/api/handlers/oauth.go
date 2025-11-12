package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"knowstack/internal/core/services"
	"knowstack/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OAuthHandler struct {
	OAuthService *services.OAuthService
}

func NewOAuthHandler(oauthService *services.OAuthService) *OAuthHandler {
	return &OAuthHandler{OAuthService: oauthService}
}

// @Summary Google Login
// @Description Redirects to Google OAuth login page
// @Tags OAuth
// @Accept json
// @Produce json
// @Success 307 {string} string "Redirect to Google OAuth login page"
// @Router /oauth/google/login [get]
func (h *OAuthHandler) GoogleLogin(c *gin.Context) {
	b := make([]byte, 32)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	c.SetCookie("oauth_state", state, 3600, "/", "", false, true)

	url := h.OAuthService.GetGoogleLoginURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// @Summary Google Callback
// @Description Handles Google OAuth callback
// @Tags OAuth
// @Accept json
// @Produce json
// @Success 200 {object} dto.GoogleAuthResponse
// @Router /oauth/google/callback [get]
func (h *OAuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	frontendURL := utils.GetEnv("FRONTEND_URL", "http://localhost:3000")

	savedState, err := c.Cookie("oauth_state")
	if err != nil || savedState != state {
		errorURL := fmt.Sprintf("%s/auth/error?message=%s", frontendURL, "Invalid state")
		c.Redirect(http.StatusTemporaryRedirect, errorURL)
		return
	}

	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	if code == "" {
		errorURL := fmt.Sprintf("%s/auth/error?message=%s", frontendURL, "Invalid code")
		c.Redirect(http.StatusTemporaryRedirect, errorURL)
		return
	}

	response, err := h.OAuthService.HandleGoogleCallback(code)
	if err != nil {
		errorURL := fmt.Sprintf("%s/auth/error?message=%s", frontendURL, "Failed to handle Google callback")
		c.Redirect(http.StatusTemporaryRedirect, errorURL)
		return
	}

	redirectURL := fmt.Sprintf("%s/oauth/google/callback#access_token=%s&refresh_token=%s&isNewUser=%t",
		frontendURL,
		response.AccessToken,
		response.RefreshToken,
		response.IsNewUser,
	)

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
