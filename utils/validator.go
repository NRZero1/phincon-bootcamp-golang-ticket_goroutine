package utils

import "github.com/go-playground/validator/v10"

var Validate *validator.Validate

func ValidateStruct(entity interface{}) error {
	Validate = validator.New(validator.WithRequiredStructEnabled())

	err := Validate.Struct(entity)

	if err != nil {
		return err	
	}
	
	return nil
}