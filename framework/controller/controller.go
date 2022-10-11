package controller

import (
	"context"
	"evoli.dev/framework/router"
	"evoli.dev/framework/validation"
)

type Base struct {
	Request router.Request
	Context context.Context
}

func (b Base) Validate(rules map[string]interface{}) validation.Validator {
	return b.Request.Validate(rules)
}

func (b Base) ValidateStruct(s interface{}) validation.Validator {
	return b.Request.ValidateStruct(s)
}
