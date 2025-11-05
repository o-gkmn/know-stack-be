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
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) CreateUser(req dto.CreateUserRequest) (*dto.CreateUserResponse, error) {

	if err := s.DB.Where("username = ?", req.Username).First(&models.User{}).Error; err == nil {
		return nil, ErrUsernameAlreadyExists
	}

	if err := s.DB.Where("email = ?", req.Email).First(&models.User{}).Error; err == nil {
		return nil, ErrEmailAlreadyExists
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := s.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return &dto.CreateUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *UserService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	var user models.User
	if err := s.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, ErrUserNotFound
	}

	if !utils.VerifyPassword(req.Password, user.Password) {
		return nil, ErrInvalidPassword
	}

	userID := strconv.FormatUint(uint64(user.ID), 10)

	token, err := utils.GenerateJWT(userID, user.Email, user.Username)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
	}, nil
}
