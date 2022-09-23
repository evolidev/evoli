package validation

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	result          map[string]interface{}
	validationError error
}

func (v Validator) Validate(data map[string]any, rules map[string]interface{}) Validator {
	val := validator.New()

	v.result = val.ValidateMap(data, rules)

	return v
}

func (v Validator) ValidateStruct(data interface{}) Validator {
	val := validator.New()

	v.validationError = val.Struct(data)
	errors := make(map[string]interface{})
	if nil != v.validationError {
		for _, err := range v.validationError.(validator.ValidationErrors) {
			errors[err.Field()] = err
		}
	}

	v.result = errors

	return v
}

func (v Validator) Errors() map[string]interface{} {
	return v.result
}

func (v Validator) Fails() bool {
	return len(v.Errors()) > 0
}

func (v Validator) Panic() {
	if len(v.Errors()) > 0 {
		panic(Error{})
	}

	if v.validationError != nil {
		panic(Error{})
	}
}

func NewValidator() Validator {
	return Validator{}
}
