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

type myTestResponseStruct struct {
	Test string
}
