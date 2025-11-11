package services

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"knowstack/internal/api/dto"
	"knowstack/internal/data/models"
	"knowstack/internal/utils"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

var (
	ErrExchangeCode  = errors.New("failed to exchange code")
	ErrGetUserInfo   = errors.New("failed to get user info")
	ErrReadResponse  = errors.New("failed to read response")
	ErrParseUserInfo = errors.New("failed to parse user info")
)

type OAuthService struct {
	DB     *gorm.DB
	config *oauth2.Config
}

func NewOAuthService(db *gorm.DB, config *oauth2.Config) *OAuthService {
	return &OAuthService{DB: db, config: config}
}

func (s *OAuthService) GetGoogleLoginURL(state string) string {
	return s.config.AuthCodeURL(state)
}

func (s *OAuthService) HandleGoogleCallback(code string) (*dto.GoogleAuthResponse, error) {
	token, err := s.config.Exchange(context.Background(), code)
	if err != nil {
		utils.LogErrorWithErr("Failed to exchange code", err)
		return nil, ErrExchangeCode
	}

	userInfo, err := s.getGoogleUserInfo(token.AccessToken)
	if err != nil {
		return nil, err
	}

	var user *models.User
	isNewUser := false

	err = s.DB.
		Preload("Role").
		Preload("Role.Claims").
		Preload("Claims").
		Where("google_id = ?", userInfo.ID).
		Or("email = ?", userInfo.Email).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		user, err = s.createGoogleUser(userInfo)
		if err != nil {
			utils.LogErrorWithErr("Failed to creating user from google", err)
			return nil, err
		}
		isNewUser = true
	} else if err != nil {
		utils.LogErrorWithErr("Failed to find user", err)
		return nil, err
	} else {
		if user.GoogleID == "" {
			user.GoogleID = userInfo.ID
			user.ProfileImage = userInfo.Picture
			user.Provider = "google"
			if err := s.DB.Save(&user).Error; err != nil {
				utils.LogErrorWithErr("Failed to update user", err)
			}
		}
	}

	userID := strconv.FormatUint(uint64(user.ID), 10)

	claimNameSet := make(map[string]struct{})
	for _, c := range user.Role.Claims {
		claimNameSet[c.Name] = struct{}{}
	}
	for _, c := range user.Claims {
		claimNameSet[c.Name] = struct{}{}
	}
	mergedClaimNames := make([]string, 0, len(claimNameSet))
	for name := range claimNameSet {
		mergedClaimNames = append(mergedClaimNames, name)
	}

	accessToken, err := utils.GenerateAccessToken(userID, user.Email, user.Username, user.RoleID, mergedClaimNames)
	if err != nil {
		utils.LogErrorWithErr("Failed to generate access token", err)
		return nil, err
	}

	refreshTokenRecord := models.RefreshToken{
		UserID:    user.ID,
		IsRevoked: false,
	}
	if err := s.DB.Create(&refreshTokenRecord).Error; err != nil {
		utils.LogErrorWithErr("Failed to create token record", err)
		return nil, err
	}

	tokenID := strconv.FormatUint(uint64(refreshTokenRecord.ID), 10)
	refreshToken, err := utils.GenerateRefreshToken(userID, tokenID, true)
	if err != nil {
		utils.LogErrorWithErr("Failed to generate refresh token", err)
		return nil, err
	}

	refreshTokenRecord.Token = refreshToken
	if err := s.DB.Save(&refreshTokenRecord).Error; err != nil {
		utils.LogErrorWithErr("Failed to save token record", err)
		return nil, err
	}

	return &dto.GoogleAuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IsNewUser:    isNewUser,
	}, nil
}

func (s *OAuthService) getGoogleUserInfo(accessToken string) (*dto.GoogleUserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		utils.LogErrorWithErr("Failed to get user info", err)
		return nil, ErrGetUserInfo
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		utils.LogError("Failed to get user info: status code %d", resp.StatusCode)
		return nil, ErrGetUserInfo
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.LogErrorWithErr("Failed to read response body", err)
		return nil, ErrReadResponse
	}

	var userInfo dto.GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		utils.LogErrorWithErr("Failed to unmarshal user info", err)
		return nil, ErrParseUserInfo
	}

	return &userInfo, nil
}

func (s *OAuthService) createGoogleUser(userInfo *dto.GoogleUserInfo) (*models.User, error) {
	var defaultRole models.Role
	if err := s.DB.Where("is_default = ?", true).First(&defaultRole).Error; err != nil {
		utils.LogErrorWithErr("Failed to find default role", err)
		return nil, ErrDefaultRoleNotFound
	}

	username := strings.Split(userInfo.Email, "@")[0]
	originalUsername := username
	counter := 1

	for {
		var existingUser models.User
		err := s.DB.Where("username = ?", username).First(&existingUser).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			break
		}
		username = originalUsername + strconv.Itoa(counter)
		counter++
	}

	user := models.User{
		Username:     username,
		Email:        userInfo.Email,
		GoogleID:     userInfo.ID,
		Provider:     "google",
		ProfileImage: userInfo.Picture,
		RoleID:       defaultRole.ID,
		Password:     "",
	}

	if err := s.DB.Create(&user).Error; err != nil {
		utils.LogErrorWithErr("Failed to create user", err)
		return nil, err
	}

	if err := s.DB.
		Preload("Role").
		Preload("Role.Claims").
		Preload("Claims").
		First(&user, user.ID).Error; err != nil {
		utils.LogErrorWithErr("Failed to get user", err)
		return nil, err
	}

	return &user, nil
}
