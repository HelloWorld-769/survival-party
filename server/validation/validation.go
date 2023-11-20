package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func CheckValidation(data interface{}) error {

	fmt.Println("data", data)
	validationErr := Validate.Struct(data)
	if validationErr != nil {
		return validationErr
	}
	return nil
}
