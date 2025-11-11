package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"knowstack/internal/core/services"
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

	savedState, err := c.Cookie("oauth_state")
	if err != nil || savedState != state {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	if code == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	response, err := h.OAuthService.HandleGoogleCallback(code)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, response)
}