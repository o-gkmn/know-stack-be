package handlers

import (
	"errors"
	"knowstack/internal/api/dto"
	"knowstack/internal/api/httperrors"
	"knowstack/internal/api/validation"
	"knowstack/internal/core/services"
	"knowstack/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

// @Summary Create a new user
// @Description Creates a new user
// @Tags API User
// @Accept json
// @Produce json
// @Success 201 {object} dto.CreateUserResponse
// @Router /users/register [post]
// @Param user body dto.CreateUserRequest true "User to create"
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if ok := utils.BindJSONAndValidate(c, &req, validation.CreateUserValidationMessages()); !ok {
		return
	}

	user, err := h.UserService.CreateUser(req)
	if err != nil {
		if errors.Is(err, services.ErrUsernameAlreadyExists) {
			httperrors.ErrUsernameAlreadyExists.Write(c)
		} else if errors.Is(err, services.ErrEmailAlreadyExists) {
			httperrors.ErrEmailAlreadyExists.Write(c)
		} else {
			httperrors.ErrInternalServerError.Write(c)
		}
		return
	}
	c.JSON(http.StatusCreated, user)
}

// @Summary Login a user
// @Description Logs in a user
// @Tags API User
// @Accept json
// @Produce json
// @Success 200 {object} dto.LoginResponse
// @Router /users/login [post]
// @Param user body dto.LoginRequest true "User to login"
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if ok := utils.BindJSONAndValidate(c, &req, validation.LoginValidationMessages()); !ok {
		return
	}
	user, err := h.UserService.Login(req)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			httperrors.ErrUserNotFound.Write(c)
		} else if errors.Is(err, services.ErrInvalidPassword) {
			httperrors.ErrInvalidPassword.Write(c)
		} else {
			httperrors.ErrInternalServerError.Write(c)
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Summary Refresh a token
// @Description Refreshes a token
// @Tags API User
// @Accept json
// @Produce json
// @Success 200 {object} dto.RefreshResponse
// @Router /users/refresh [post]
// @Param user body dto.RefreshRequest true "User to refresh"
func (h *UserHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if ok := utils.BindJSONAndValidate(c, &req, validation.RefreshValidationMessages()); !ok {
		return
	}
	res, err := h.UserService.Refresh(req)
	if err != nil {
		if errors.Is(err, utils.ErrTokenExpired) {
			httperrors.ErrTokenExpired.Write(c)
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
	c.JSON(http.StatusOK, res)
}

// @Summary Logout a user
// @Description Logs out a user
// @Tags API User
// @Accept json
// @Produce json
// @Success 200 {object} dto.LogoutResponse
// @Router /users/logout [post]
// @Param user body dto.LogoutRequest true "User to logout"
func (h *UserHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if ok := utils.BindJSONAndValidate(c, &req, validation.LogoutValidationMessages()); !ok {
		return
	}
	res, err := h.UserService.Logout(req)
	if err != nil {
		httperrors.ErrInternalServerError.Write(c)
	}
	c.JSON(http.StatusOK, res)
}

// @Summary Set claims for a user
// @Description Sets claims for a user
// @Tags API User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.SetClaimsResponse
// @Router /users/claims [post]
// @Param user body dto.SetClaimsRequest true "User to set claims for"
func (h *UserHandler) SetClaims(c *gin.Context) {
	var req dto.SetClaimsRequest
	if ok := utils.BindJSONAndValidate(c, &req, validation.SetClaimsValidationMessages()); !ok {
		return
	}
	err := h.UserService.SetClaims(req.UserID, req.ClaimIDs)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			httperrors.ErrUserNotFound.Write(c)
		} else if errors.Is(err, services.ErrClaimsNotFound) {
			httperrors.ErrClaimsNotFound.Write(c)
		} else {
			httperrors.ErrInternalServerError.Write(c)
		}
		return
	}

	c.Status(http.StatusOK)
}
