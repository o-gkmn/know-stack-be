package handlers

import (
	"errors"
	"knowstack/internal/api/dto"
	"knowstack/internal/api/httperrors"
	"knowstack/internal/api/validation"
	"knowstack/internal/core/service"
	"knowstack/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

// @Summary Create a new user
// @Description Creates a new user
// @Tags API User
// @Accept json
// @Produce json
// @Success 201 {object} dto.CreateUserResponse
// @Router /users [post]
// @Param user body dto.CreateUserRequest true "User to create"
// @Response 400 {object} httperrors.HTTPValidationError
// @Response 409 {object} httperrors.HTTPError
// @Response 500 {object} httperrors.HTTPError
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if ok := utils.BindJSONAndValidate(c, &req, validation.CreateUserValidationMessages()); !ok {
		return
	}

	user, err := h.UserService.CreateUser(req)
	if err != nil {
		if errors.Is(err, service.ErrUsernameAlreadyExists) {
			httperrors.ErrUsernameAlreadyExists.Write(c)
		} else if errors.Is(err, service.ErrEmailAlreadyExists) {
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
// @Response 400 {object} httperrors.HTTPValidationError
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if ok := utils.BindJSONAndValidate(c, &req, validation.LoginValidationMessages()); !ok {
		return
	}
	user, err := h.UserService.Login(req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			httperrors.ErrUserNotFound.Write(c)
		} else if errors.Is(err, service.ErrInvalidPassword) {
			httperrors.ErrInvalidPassword.Write(c)
		} else {
			httperrors.ErrInternalServerError.Write(c)
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Summary Set claims for a user
// @Description Sets claims for a user
// @Tags API User
// @Accept json
// @Produce json
// @Success 200 {object} dto.SetClaimsResponse
// @Router /users/claims [post]
// @Param user body dto.SetClaimsRequest true "User to set claims for"
// @Response 400 {object} httperrors.HTTPValidationError
func (h *UserHandler) SetClaims(c *gin.Context) {
	var req dto.SetClaimsRequest
	if ok := utils.BindJSONAndValidate(c, &req, validation.SetClaimsValidationMessages()); !ok {
		return
	}
	err := h.UserService.SetClaims(req.UserID, req.ClaimIDs)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			httperrors.ErrUserNotFound.Write(c)
		} else if errors.Is(err, service.ErrClaimsNotFound) {
			httperrors.ErrClaimsNotFound.Write(c)
		} else {
			httperrors.ErrInternalServerError.Write(c)
		}
		return
	}

	c.Status(http.StatusOK)
}
