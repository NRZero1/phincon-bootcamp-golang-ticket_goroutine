package utils

import "github.com/go-playground/validator/v10"

func ValidateStruct(entity interface{}) error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(entity)

	if err != nil {
		return err	
	}
	
	return nil
}