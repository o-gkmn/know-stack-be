package httperrors

import "github.com/gin-gonic/gin"

type HTTPError struct {
	Code  int    `json:"code"`
	Type  string `json:"type"`
	Title string `json:"title"`
}

type HTTPValidationError struct {
	HTTPError       `json:"http_error"`
	ValidationError ValidationErrors `json:"validation_error"`
}

type ValidationErrors struct {
	Error string `json:"error"`
	Key   string `json:"key"`
	In    string `json:"in"`
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
