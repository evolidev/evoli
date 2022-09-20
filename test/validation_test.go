package test

import (
	evoli "github.com/evolidev/evoli/framework/router"
	"github.com/evolidev/evoli/framework/validation"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

func TestName(t *testing.T) {
	t.Run("Validate should return no errors", func(t *testing.T) {
		data := make(map[string]any)
		data["test"] = "hello"

		rules := make(map[string]interface{})
		rules["test"] = "lowercase"

		v := validation.NewValidator()

		result := v.Validate(data, rules)

		assert.Empty(t, result)
	})
}

func TestRequestValidation(t *testing.T) {
	t.Run("Request should validate", func(t *testing.T) {
		r := evoli.NewRouter()
		f := url.Values{}
		f.Set("form", "myform")
		r.Post("/request/form/validation", func(request *evoli.Request) string {
			rules := make(map[string]interface{})
			rules["form"] = "lowercase"
			request.Validate(rules)

			return "pass"
		})

		rr := sendRequestWithForm(t, r, http.MethodGet, "/request/form/validation", f)

		assert.Equal(t, "pass", rr.Body.String())
	})

	t.Run("Request should validate input struct", func(t *testing.T) {
		r := evoli.NewRouter()
		f := url.Values{}
		f.Set("test", "myform")
		r.Post("/request/form/validation", func(request *evoli.Request, inputStruct InputWithValidationRules) string {
			request.ValidateStruct(inputStruct)

			return "pass"
		})

		rr := sendRequestWithForm(t, r, http.MethodGet, "/request/form/validation", f)

		assert.Equal(t, "pass", rr.Body.String())
	})
}

type InputWithValidationRules struct {
	Test string `validate:"lowercase"`
}
