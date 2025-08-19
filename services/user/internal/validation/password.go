package validation

import "github.com/go-playground/validator"

func RegisterCustomValidations(v *validator.Validate) {
	v.RegisterValidation("passwd", validatePassword)
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	hasLetter := false
	hasDigit := false

	for _, ch := range password {
		switch {
		case 'a' <= ch && ch <= 'z', 'A' <= ch && ch <= 'Z':
			hasLetter = true
		case '0' <= ch && ch <= '9':
			hasDigit = true
		}
	}

	return hasLetter && hasDigit
}
