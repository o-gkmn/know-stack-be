package httperrors

import "github.com/gin-gonic/gin"

type HTTPError struct {
	Code  int
	Type  string
	Title string
}

type HTTPValidationError struct {
	HTTPError
	ValidationError ValidationErrors
}

type ValidationErrors struct {
	Error string
	Key   string
	In    string
}

func NewHTTPError(code int, errorType, title string) *HTTPError {
	return &HTTPError{
		Code:  code,
		Type:  errorType,
		Title: title,
	}
}

func NewHTTPValidationError(code int, errorType, title string, validationError ValidationErrors) *HTTPValidationError {
	return &HTTPValidationError{
		HTTPError: HTTPError{
			Code:  code,
			Type:  errorType,
			Title: title,
		},
		ValidationError: validationError,
	}
}

func (h *HTTPError) Write(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(h.Code, h)
}

func (h *HTTPValidationError) Write(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(h.Code, h)
}
