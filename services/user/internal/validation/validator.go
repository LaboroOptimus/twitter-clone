package validation

import "github.com/go-playground/validator"

var Validate *validator.Validate

func Init() {
	Validate = validator.New()
	RegisterCustomValidations(Validate)
}
