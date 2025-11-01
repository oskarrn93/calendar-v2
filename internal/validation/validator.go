package validator

import (
	govalidator "github.com/go-playground/validator/v10"
)

var validate = govalidator.New(govalidator.WithPrivateFieldValidation())

func ValidateStruct(val any) error {
	err := validate.Struct(val)
	if err != nil {
		return err
	}
	return nil
}
