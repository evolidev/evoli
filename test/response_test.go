package test

import (
	"encoding/json"
	evoli "github.com/evolidev/evoli/framework/response"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResponseToByteArray(t *testing.T) {
	t.Run("StringResponse should be able to get data as string", func(t *testing.T) {
		resp := evoli.String("hello-world")

		assert.Exactly(t, "hello-world", string(resp.AsBytes()))
	})

	t.Run("JsonResponse should be able to get data as string", func(t *testing.T) {
		myStruct := myTestResponseStruct{Test: "test"}
		resp := evoli.Json(myStruct)

		result, err := json.Marshal(myStruct)
		if err != nil {
			t.Fatal(err)
		}
		assert.Exactly(t, string(result), string(resp.AsBytes()))
	})
}

func TestFactory(t *testing.T) {
	t.Run("Factory should return StringResponse if input is already a StringResponse", func(t *testing.T) {
		resp := evoli.NewResponse(evoli.String("hello-world"))

		assert.Exactly(t, "hello-world", string(resp.AsBytes()))
	})

	t.Run("Factory should return StringResponse if input is a string", func(t *testing.T) {
		resp := evoli.NewResponse("hello-world")

		assert.Exactly(t, "hello-world", string(resp.AsBytes()))
	})

	t.Run("Factory should return JsonResponse if input is a struct", func(t *testing.T) {
		myStruct := myTestResponseStruct{Test: "test"}
		resp := evoli.NewResponse(myStruct)
		result, err := json.Marshal(myStruct)
		if err != nil {
			t.Fatal(err)
		}

		assert.Exactly(t, string(result), string(resp.AsBytes()))
	})
}

func TestView(t *testing.T) {
	t.Run("ViewResponse should return value from template", func(t *testing.T) {
		view := evoli.View("templates.test")

		assert.Exactly(t, "<div>Hello test</div>", string(view.AsBytes()))
	})

	t.Run("View withHeader should add given key value to headers list", func(t *testing.T) {
		viewResponse := evoli.View("test")

		viewResponse.WithHeader("My-Header", "Test")

		assert.Exactly(t, "Test", viewResponse.Headers().Get("My-Header"))
	})

	t.Run("ViewResponse should handle given data", func(t *testing.T) {
		data := map[string]any{"Name": "test"}

		view := evoli.View("templates.test_with_data").
			WithData(data)

		assert.Exactly(t, "<div>Hello test</div>", string(view.AsBytes()))
	})

	t.Run("ViewResponse should parse layout", func(t *testing.T) {
		view := evoli.View("templates.test_with_layout").
			WithLayout("templates.layout")

		assert.Exactly(t, "<main><div>Hello test layout</div></main>", string(view.AsBytes()))
	})
}

type myTestResponseStruct struct {
	Test string
}
