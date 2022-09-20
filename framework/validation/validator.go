package validation

import "github.com/go-playground/validator/v10"

type Validator struct {
}

func (v Validator) Validate(data map[string]any, rules map[string]interface{}) map[string]interface{} {
	val := validator.New()

	return val.ValidateMap(data, rules)
}

func (v Validator) ValidateStruct(data interface{}) error {
	val := validator.New()

	return val.Struct(data)
}

func NewValidator() Validator {
	return Validator{}
}
