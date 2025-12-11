package util

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

const indonesianPhoneRegex = `^(\+62|62|0)(\d{8,15})$`

func RegisterCustomValidators(v *validator.Validate) {
	err := v.RegisterValidation("phone_id", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		if phone == "" {
			return true
		}
		match, _ := regexp.MatchString(indonesianPhoneRegex, phone)
		return match
	})
	if err != nil {
		panic(err)
	}
}
