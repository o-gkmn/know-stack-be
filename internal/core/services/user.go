package services

import (
	"errors"
	"fmt"
	"knowstack/internal/api/dto"
	"knowstack/internal/data/models"
	"knowstack/internal/utils"
	"strconv"
	"time"

	"gorm.io/gorm"
)

var (
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrUserNotFound          = errors.New("user not found")
	ErrDefaultRoleNotFound   = errors.New("default role not found")
	ErrClaimsNotFound        = errors.New("claims not found")
	ErrTokenNotFound         = errors.New("token not found")
	ErrInvalidToken          = errors.New("invalid token")
	ErrParseError            = errors.New("parsing error")
	ErrMismatchTokenAndUser  = errors.New("token and user mismatch")
	ErrRefreshTokenExpired   = errors.New("refresh token expired")
)

type UserService struct {
	DB    *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		DB:    db,
	}
}

func (s *UserService) CreateUser(req dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	utils.LogInfo("Creating user: %+v", req)

	if err := s.DB.Where("username = ?", req.Username).First(&models.User{}).Error; err == nil {
		utils.LogInfo("Username already exists: %+v", req.Username)
		return nil, ErrUsernameAlreadyExists
	}

	if err := s.DB.Where("email = ?", req.Email).First(&models.User{}).Error; err == nil {
		utils.LogInfo("Email already exists: %+v", req.Email)
		return nil, ErrEmailAlreadyExists
	}

	var defaultRole models.Role
	if err := s.DB.Where("is_default = ?", true).First(&defaultRole).Error; err != nil {
		utils.LogErrorWithErr("Failed to find default role", err)
		return nil, ErrDefaultRoleNotFound
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		RoleID:   defaultRole.ID,
		Password: req.Password,
	}

	if err := s.DB.Create(user).Error; err != nil {
		utils.LogErrorWithErr("Failed to create user", err)
		return nil, err
	}

	return &dto.CreateUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *UserService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	utils.LogInfo("Logging in user: %+v", req)

	var user models.User
	if err := s.DB.
		Preload("Role").
		Preload("Role.Claims").
		Preload("Claims").
		Where("email = ?", req.Email).
		First(&user).Error; err != nil {
		utils.LogErrorWithErr("Failed to find user", err)
		return nil, ErrUserNotFound
	}

	if !utils.VerifyPassword(req.Password, user.Password) {
		utils.LogError("Invalid password")
		return nil, ErrInvalidPassword
	}

	userID := strconv.FormatUint(uint64(user.ID), 10)

	// Merge role claims and user claims distinctly
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

	token, err := utils.GenerateAccessToken(userID, user.Email, user.Username, user.RoleID, mergedClaimNames)
	if err != nil {
		utils.LogErrorWithErr("Failed to generate JWT", err)
		return nil, err
	}

	// Initialize refresh token record with user ID; token will be assigned after ID generation
	refreshTokenRecord := models.RefreshToken{
		UserID:    user.ID,
		IsRevoked: false,
	}
	if err := s.DB.Create(&refreshTokenRecord).Error; err != nil {
		utils.LogErrorWithErr("Failed to create token record", err)
		return nil, err
	}

	tokenID := strconv.FormatUint(uint64(refreshTokenRecord.ID), 10)
	refreshToken, err := utils.GenerateRefreshToken(userID, tokenID, req.Remember)
	if err != nil {
		utils.LogErrorWithErr("Failed to generate refresh token", err)
		return nil, err
	}

	refreshTokenRecord.Token = refreshToken
	if err := s.DB.Save(&refreshTokenRecord).Error; err != nil {
		utils.LogErrorWithErr("Failed to save token record", err)
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) Refresh(req dto.RefreshRequest) (*dto.RefreshResponse, error) {
	utils.LogInfo("Refreshing token: %+v", req)

	claims, err := utils.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		utils.LogErrorWithErr("Failed to validate refresh token", err)
		return nil, err
	}

	// Parse token ID from claims
	tokenID64, err := strconv.ParseUint(claims.TokenID, 10, 32)
	if err != nil {
		utils.LogErrorWithErr("Failed to parse token ID", err)
		return nil, err
	}
	tokenID := uint(tokenID64)

	// Parse user ID from claims
	userID64, err := strconv.ParseUint(claims.UserID, 10, 32)
	if err != nil {
		utils.LogErrorWithErr("Failed to parse user ID", err)
		return nil, err
	}
	userID := uint(userID64)

	var token models.RefreshToken
	if err := s.DB.
		Preload("User").
		Preload("User.Claims").
		Preload("User.Role.Claims").
		Where("ID = ?", tokenID).
		First(&token).Error; err != nil {
		utils.LogErrorWithErr("Failed to find token", err)
		return nil, ErrTokenNotFound
	}

	if userID != token.UserID {
		utils.LogError("Token and user mismatch")
		return nil, ErrMismatchTokenAndUser
	}

	// merged role and user claims
	claimNameSet := make(map[string]struct{})
	for _, c := range token.User.Role.Claims {
		claimNameSet[c.Name] = struct{}{}
	}
	for _, c := range token.User.Claims {
		claimNameSet[c.Name] = struct{}{}
	}
	mergedClaimNames := make([]string, 9, len(claimNameSet))
	for name := range claimNameSet {
		mergedClaimNames = append(mergedClaimNames, name)
	}

	accessToken, err := utils.GenerateAccessToken(
		claims.UserID,
		token.User.Email,
		token.User.Username,
		token.User.RoleID,
		mergedClaimNames,
	)
	if err != nil {
		utils.LogErrorWithErr("Failed to generate access token", err)
		return nil, err
	}

	return &dto.RefreshResponse{
		AccessToken: accessToken,
	}, nil
}

func (s *UserService) RequestPasswordReset(req dto.RequestPasswordResetRequest) (*dto.RequestPasswordResetResponse, error) {
	utils.LogInfo("Requesting password reset for email: %+v", req.Email)

	var user models.User
	if err := s.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogInfo("User not found: %+v", req.Email)
			// Don't reveal if user exists or not for security reasons
			// Return success even if user doesn't exist
			return &dto.RequestPasswordResetResponse{IsSuccess: true}, nil
		}

		utils.LogErrorWithErr("Failed to find user", err)
		return nil, err
	}

	// Generate password reset token
	token, err := utils.GeneratePasswordResetToken()
	if err != nil {
		utils.LogErrorWithErr("Failed to generate password reset token", err)
		return nil, err
	}

	// Set expiration time (default: 1 hour)
	expiresInHours := utils.GetEnvAsInt("PASSWORD_RESET_EXPIRES_IN_HOURS", 1)
	expiresAt := time.Now().Add(time.Duration(expiresInHours) * time.Hour)

	// Revoke any existing unused tokens for this user
	if err := s.DB.Model(&models.PasswordResetToken{}).
		Where("user_id = ? AND is_used = ?", user.ID, false).
		Update("is_used", true).Error; err != nil {
		utils.LogErrorWithErr("Failed to revoke existing tokens", err)
		// Continue anyway, not critical
	}

	// Create new password reset token
	passwordResetToken := models.PasswordResetToken{
		Token:     token,
		ExpiresAt: expiresAt,
		IsUsed:    false,
		UserID:    user.ID,
	}

	if err := s.DB.Create(&passwordResetToken).Error; err != nil {
		utils.LogErrorWithErr("Failed to create password reset token", err)
		return nil, err
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", utils.GetEnv("FRONTEND_URL", "http://localhost:3000"), token)
	// TODO: Send email with reset URL
	utils.LogInfo("Reset URL: %s", resetURL)
	utils.LogInfo("Password reset email sent", "email", user.Email)

	return &dto.RequestPasswordResetResponse{IsSuccess: true}, nil
}

func (s *UserService) Logout(req dto.LogoutRequest) (*dto.LogoutResponse, error) {
	if err := s.DB.Model(&models.RefreshToken{}).
		Where("token = ?", req.RefreshToken).
		Update("is_revoked", true).Error; err != nil {
		utils.LogErrorWithErr("Failed to logout", err)
		return &dto.LogoutResponse{IsSuccess: false}, err
	}

	return &dto.LogoutResponse{IsSuccess: true}, nil
}

// SetClaims adds new claims and removes existing claims from the user
func (s *UserService) SetClaims(userID uint, claimIDs []uint) error {
	utils.LogInfo("Setting claims for user: %+v", userID)

	var user models.User
	if err := s.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		utils.LogErrorWithErr("Failed to find user", err)
		return ErrUserNotFound
	}

	var claims []models.Claim
	if err := s.DB.Where("id IN (?)", claimIDs).Find(&claims).Error; err != nil {
		utils.LogErrorWithErr("Failed to find claims", err)
		return ErrClaimsNotFound
	}

	user.Claims = claims

	if err := s.DB.Save(&user).Error; err != nil {
		utils.LogErrorWithErr("Failed to save user", err)
		return err
	}

	return nil
}
