package utils

import (
	"errors"
	"net/http"

	"knowstack/internal/api/httperrors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// FieldErrorMessages allows per-field, per-tag custom messages.
// Example: {"Username": {"required": "...", "min": "..."}}
type FieldErrorMessages map[string]map[string]string

// BindJSONAndValidate binds JSON into dst and, on error, writes a single validation error response.
// Returns true if binding/validation succeeded; false if an error response was written.
func BindJSONAndValidate(c *gin.Context, dst any, messages FieldErrorMessages) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		var verrs validator.ValidationErrors
		msg := err.Error()
		key := "request"
		in := "body"
		if errors.As(err, &verrs) && len(verrs) > 0 {
			fe := verrs[0]
			key = fe.Field()
			if fieldMsgs, ok := messages[key]; ok {
				if m, found := fieldMsgs[fe.Tag()]; found {
					msg = m
				} else {
					msg = defaultReadableMessage(fe)
				}
			} else {
				msg = defaultReadableMessage(fe)
			}
		}

		valErr := httperrors.NewHTTPValidationError(
			http.StatusBadRequest,
			"validation_error",
			"Validation error",
			httperrors.ValidationErrors{
				Error: msg,
				Key:   key,
				In:    in,
			},
		)
		valErr.Write(c)
		return false
	}
	return true
}

func defaultReadableMessage(fe validator.FieldError) string {
	// Generic fallback if no custom message provided
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fe.Field() + " is too short"
	case "max":
		return fe.Field() + " is too long"
	case "alphanumunicode":
		return fe.Field() + " must be alphanumeric"
	default:
		return fe.Error()
	}
}
