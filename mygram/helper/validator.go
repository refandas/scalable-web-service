package helper

import "github.com/go-playground/validator/v10"

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) Validate(s any) error {
	return v.validate.Struct(s)
}
