package validation

import "knowstack/internal/utils"

// CreateUserValidationMessages returns field-specific, tag-specific messages for CreateUserRequest.
func CreateUserValidationMessages() utils.FieldErrorMessages {
	return utils.FieldErrorMessages{
		"Username": {
			"required":        "Kullanıcı adı zorunludur.",
			"alphanumunicode": "Kullanıcı adı yalnızca harf ve rakam içerebilir.",
			"min":             "Kullanıcı adı en az 3 karakter olmalıdır.",
			"max":             "Kullanıcı adı en fazla 30 karakter olabilir.",
		},
		"Email": {
			"required": "E-posta zorunludur.",
			"email":    "Geçerli bir e-posta adresi giriniz.",
		},
		"Password": {
			"required": "Parola zorunludur.",
			"min":      "Parola en az 8 karakter olmalıdır.",
			"max":      "Parola en fazla 72 karakter olabilir.",
		},
	}
}

// LoginValidationMessages returns field-specific, tag-specific messages for LoginRequest.
func LoginValidationMessages() utils.FieldErrorMessages {
	return utils.FieldErrorMessages{
		"Email": {
			"required": "E-posta zorunludur.",
			"email":    "Geçerli bir e-posta adresi giriniz.",
		},
		"Password": {
			"required": "Parola zorunludur.",
			"min":      "Parola en az 8 karakter olmalıdır.",
			"max":      "Parola en fazla 72 karakter olabilir.",
		},
		"Remember": {
			"required": "Remember me zorunludur",
			"boolean":  "Remember me bir boolean olmalıdır",
		},
	}
}

func RefreshValidationMessages() utils.FieldErrorMessages {
	return utils.FieldErrorMessages{
		"RefreshToken": {
			"required": "Refresh token zorunludur",
		},
	}
}

func LogoutValidationMessages() utils.FieldErrorMessages {
	return utils.FieldErrorMessages{
		"RefreshToken": {
			"required": "Refresh token zorunludur",
		},
	}
}

func RequestPasswordResetValidationMessages() utils.FieldErrorMessages {
	return utils.FieldErrorMessages{
		"Email": {
			"required": "E-posta zorunludur.",
		},
	}
}

func SetClaimsValidationMessages() utils.FieldErrorMessages {
	return utils.FieldErrorMessages{
		"UserID": {
			"required": "Kullanıcı ID zorunludur.",
		},
		"ClaimIDs": {
			"required": "Claims zorunludur.",
		},
	}
}