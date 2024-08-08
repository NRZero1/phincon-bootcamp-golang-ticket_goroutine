package utils

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func ValidateStruct(entity interface{}) error {
	Validate = validator.New(validator.WithRequiredStructEnabled())

	err := Validate.Struct(entity)

	if err != nil {
		return err	
	}
	
	return nil
}

func VerifyPassword(s string) (sevenOrMore, number, upper bool) {
    letters := 0
    for _, c := range s {
        switch {
        case unicode.IsNumber(c):
            number = true
        case unicode.IsUpper(c):
            upper = true
            letters++
        case unicode.IsLetter(c) || c == ' ':
            letters++
		}
    }
    sevenOrMore = letters >= 7
    return
}