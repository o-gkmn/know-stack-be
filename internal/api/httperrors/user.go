package httperrors

import (
	"net/http"
)

var (
	ErrInvalidRequest        = NewHTTPError(http.StatusBadRequest, "invalid_request", "Invalid request")
	ErrInternalServerError   = NewHTTPError(http.StatusInternalServerError, "internal_server_error", "Internal server error")
	ErrUsernameAlreadyExists = NewHTTPError(http.StatusConflict, "username_already_exists", "Username already exists")
	ErrEmailAlreadyExists    = NewHTTPError(http.StatusConflict, "email_already_exists", "Email already exists")
	ErrUserNotFound          = NewHTTPError(http.StatusNotFound, "user_not_found", "User not found")
	ErrInvalidPassword       = NewHTTPError(http.StatusBadRequest, "invalid_password", "Invalid password")
)
