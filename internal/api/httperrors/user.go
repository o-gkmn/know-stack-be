package httperrors

import (
	"net/http"
)

var (
	ErrInvalidRequest        = NewHTTPError(http.StatusBadRequest, "invalid_request", "Geçersiz istek")
	ErrInternalServerError   = NewHTTPError(http.StatusInternalServerError, "internal_server_error", "Sunucu hatası")
	ErrUsernameAlreadyExists = NewHTTPError(http.StatusConflict, "username_already_exists", "Kullanıcı adı zaten kullanılıyor")
	ErrEmailAlreadyExists    = NewHTTPError(http.StatusConflict, "email_already_exists", "E-posta zaten kullanılıyor")
	ErrUserNotFound          = NewHTTPError(http.StatusNotFound, "user_not_found", "Kullanıcı bulunamadı")
	ErrInvalidPassword       = NewHTTPError(http.StatusBadRequest, "invalid_password", "Geçersiz şifre")
	ErrClaimsNotFound        = NewHTTPError(http.StatusNotFound, "claims_not_found", "Yetkinlikler bulunamadı")
	ErrTokenExpired          = NewHTTPError(http.StatusUnauthorized, "token_expired", "Token süresi dolmuş.")
)
