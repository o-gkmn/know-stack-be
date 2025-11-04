package httperrors

import "github.com/gin-gonic/gin"

type HTTPError struct {
	Code  int
	Type  string
	Title string
}

type HTTPValidationError struct {
	HTTPError
	ValidationErrors []ValidationErrors
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

func NewHTTPValidationError(code int, errorType, title string, validationErrors []ValidationErrors) *HTTPValidationError {
	return &HTTPValidationError{
		HTTPError: HTTPError{
			Code:  code,
			Type:  errorType,
			Title: title,
		},
		ValidationErrors: validationErrors,
	}
}

func (h *HTTPError) Write(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(h.Code, h)
}

func (h *HTTPValidationError) Write(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(h.Code, h)
}
