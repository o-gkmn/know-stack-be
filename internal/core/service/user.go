package service

import (
	"errors"
	"knowstack/internal/api/dto"
	"knowstack/internal/data/models"
	"knowstack/internal/utils"
	"strconv"

	"gorm.io/gorm"
)

var (
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrUserNotFound          = errors.New("user not found")
	ErrDefaultRoleNotFound   = errors.New("default role not found")
	ErrClaimsNotFound        = errors.New("claims not found")
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
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

	token, err := utils.GenerateJWT(userID, user.Email, user.Username, user.RoleID, mergedClaimNames)
	if err != nil {
		utils.LogErrorWithErr("Failed to generate JWT", err)
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
	}, nil
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
