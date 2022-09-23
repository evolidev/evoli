package test

import (
	"github.com/evolidev/evoli/framework/controller"
	"github.com/evolidev/evoli/framework/router"
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

		assert.False(t, result.Fails())
	})
}

func TestRequestValidation(t *testing.T) {
	t.Run("Request should validate", func(t *testing.T) {
		r := router.NewRouter()
		f := url.Values{}
		f.Set("form", "myform")
		r.Post("/request/form/validation", func(request *router.Request) string {
			rules := make(map[string]interface{})
			rules["form"] = "lowercase"
			request.Validate(rules).Panic()

			return "pass"
		})

		rr := sendRequestWithForm(t, r, http.MethodGet, "/request/form/validation", f)

		assert.Equal(t, "pass", rr.Body.String())
	})

	t.Run("Request should validate input struct", func(t *testing.T) {
		r := router.NewRouter()
		f := url.Values{}
		f.Set("test", "myform")
		r.Post("/request/form/validation", func(request *router.Request, inputStruct InputWithValidationRules) string {
			request.ValidateStruct(inputStruct).Panic()

			return "pass"
		})

		rr := sendRequestWithForm(t, r, http.MethodGet, "/request/form/validation", f)

		assert.Equal(t, "pass", rr.Body.String())
	})

	t.Run("Request return error", func(t *testing.T) {
		r := router.NewRouter()
		f := url.Values{}
		f.Set("test", "Myform")
		r.Post("/request/form/validation", func(request *router.Request, inputStruct InputWithValidationRules) string {
			result := request.ValidateStruct(inputStruct)

			if result.Fails() {
				return "fail"
			}

			return "pass"
		})

		rr := sendRequestWithForm(t, r, http.MethodGet, "/request/form/validation", f)

		assert.Equal(t, "fail", rr.Body.String())
	})
}

func TestController(t *testing.T) {
	t.Run("Request should validate input struct", func(t *testing.T) {
		r := router.NewRouter()
		f := url.Values{}
		f.Set("form", "myform")
		r.Post("/request/controller/form/validation", ValidationController.TestAction)

		rr := sendRequestWithForm(t, r, http.MethodGet, "/request/form/validation", f)

		assert.Equal(t, "pass", rr.Body.String())
	})

	t.Run("Request should validate input struct", func(t *testing.T) {
		r := router.NewRouter()
		f := url.Values{}
		f.Set("test", "myform")
		r.Post("/request/controller/form/struct/validation", ValidationController.TestInputAction)

		rr := sendRequestWithForm(t, r, http.MethodGet, "/request/controller/form/struct/validation", f)

		assert.Equal(t, "pass", rr.Body.String())
	})
}

type InputWithValidationRules struct {
	Test string `validate:"lowercase"`
}

type ValidationController struct {
	controller.Base
}

func (c ValidationController) TestAction() string {
	rules := make(map[string]interface{})
	rules["form"] = "lowercase"
	c.Validate(rules).Panic()

	return "pass"
}

func (c ValidationController) TestInputAction(inputStruct InputWithValidationRules) string {
	c.ValidateStruct(inputStruct).Panic()

	return "pass"
}
