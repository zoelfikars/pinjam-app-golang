package util

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func TranslateValidationErrors(err error) map[string][]string {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}
	errorMap := make(map[string][]string)
	for _, fieldError := range validationErrors {
		fieldName := ToSnakeCase(fieldError.Field())
		message := generateErrorMessage(fieldError)
		errorMap[fieldName] = append(errorMap[fieldName], message)
	}
	return errorMap
}
func generateErrorMessage(fe validator.FieldError) string {
	fieldName := cases.Title(language.Indonesian).String(strings.ReplaceAll(ToSnakeCase(fe.Field()), "_", " "))
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s diperlukan dan tidak boleh kosong.", fieldName)
	case "phone_id":
		return fmt.Sprintf("%s harus dalam format nomor telepon Indonesia yang valid (diawali 0 atau +62).", fieldName)
	case "numeric":
		return fmt.Sprintf("%s harus berupa angka.", fieldName)
	case "email":
		return fmt.Sprintf("%s harus berupa alamat email yang valid.", fieldName)
	case "min":
		kind := fe.Type().Kind()
		if kind == reflect.String {
			return fmt.Sprintf("%s minimal harus %s karakter.", fieldName, fe.Param())
		}
		return fmt.Sprintf("%s minimal harus %s.", fieldName, fe.Param())
	case "gte":
		return fmt.Sprintf("%s harus lebih besar atau sama dengan %s.", fieldName, fe.Param())
	case "gt":
		return fmt.Sprintf("%s harus lebih besar dari %s.", fieldName, fe.Param())
	case "oneof":
		allowedValues := strings.ReplaceAll(fe.Param(), " ", ", ")
		return fmt.Sprintf("%s harus salah satu dari nilai berikut: %s.", fieldName, allowedValues)
	case "eqfield":
		return fmt.Sprintf("%s tidak cocok dengan field %s.", fieldName, fe.Param())
	case "required_if":
		params := strings.Split(fe.Param(), " ")
		if len(params) == 2 {
			syaratField := cases.Title(language.Indonesian).String(params[0])
			syaratValue := params[1]
			return fmt.Sprintf("%s diperlukan karena %s diset sebagai '%s'.", fieldName, syaratField, syaratValue)
		}
		return fmt.Sprintf("%s diperlukan berdasarkan kondisi lain.", fieldName)
	default:
		return fmt.Sprintf("%s gagal dalam validasi tag '%s'.", fieldName, fe.Tag())
	}
}
