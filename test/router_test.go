package test

import (
	"encoding/json"
	evoli "github.com/evolidev/evoli/framework"
	"github.com/evolidev/evoli/framework/response"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBasic(t *testing.T) {
	router := evoli.NewRouter()

	for method, fn := range router.MethodTable() {
		t.Run("Basic "+method+" route should have the returned string of the callback in the body", func(t *testing.T) {
			path := "/" + strings.ToLower(method)

			fn(path, handler)

			rr := sendRequest(t, router, method, path)
			assert.Exactly(t, "hello-world", rr.Body.String())
		})
	}

	t.Run("Basic match route should have the returned string in body", func(t *testing.T) {
		path := "/test"
		router.Match(path, handler, http.MethodGet, http.MethodPost)

		rr := sendRequest(t, router, http.MethodGet, path)
		assert.Exactly(t, "hello-world", rr.Body.String())

		rr = sendRequest(t, router, http.MethodPut, path)
		assert.Exactly(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("Basic any route should have the returned string in body for all http methods", func(t *testing.T) {
		path := "/any"
		router.Any(path, handler)

		for method, _ := range router.MethodTable() {
			rr := sendRequest(t, router, method, path)
			assert.Exactly(t, "hello-world", rr.Body.String())
		}
	})

	t.Run("Basic route should be able to return a struct or slice which then get converted to json", func(t *testing.T) {
		pathStruct := "/struct"
		pathSlice := "/slice"
		router.Get(pathStruct, structHandler)
		router.Get(pathSlice, sliceHandler)

		rr := sendRequest(t, router, http.MethodGet, pathStruct)

		testJson, err := json.Marshal(testStruct{"test"})
		if err != nil {
			t.Fatal(err)
		}

		assert.Exactly(t, string(testJson), rr.Body.String())
		assert.Exactly(t, "application/json", rr.Header().Get("Content-Type"))

		rr = sendRequest(t, router, http.MethodGet, pathSlice)

		testJson, err = json.Marshal([]uint8{255, 255, 255})
		if err != nil {
			t.Fatal(err)
		}

		assert.Exactly(t, string(testJson), rr.Body.String())
		assert.Exactly(t, "application/json", rr.Header().Get("Content-Type"))
	})

	t.Run("Basic route should return plain int if return is an int", func(t *testing.T) {
		path := "/int"
		router.Get(path, handlerInt)

		rr := sendRequest(t, router, http.MethodGet, path)

		assert.Exactly(t, "1", rr.Body.String())
		assert.Exactly(t, "text/plain", rr.Header().Get("Content-Type"))
	})

	t.Run("Basic route should be able to handle response object", func(t *testing.T) {
		path := "/response/string"
		router.Get(path, handlerStringResponse)

		rr := sendRequest(t, router, http.MethodGet, path)

		assert.Exactly(t, "hello-world", rr.Body.String())

		path = "/response/json"
		router.Get(path, handlerJsonResponse)

		rr = sendRequest(t, router, http.MethodGet, path)
		testJson, err := json.Marshal(testStruct{Test: "test"})
		if err != nil {
			t.Fatal(err)
		}
		assert.Exactly(t, string(testJson), rr.Body.String())
	})
}

func handler() interface{} {
	return "hello-world"
}

func handlerInt() interface{} {
	return 1
}

func handlerStringResponse() interface{} {
	return response.String("hello-world")
}

func handlerJsonResponse() interface{} {
	return response.Json(testStruct{Test: "test"})
}

// todo do not loose return type
func structHandler() interface{} {
	return testStruct{"test"}
}

func sliceHandler() interface{} {
	return []uint8{255, 255, 255}
}

type testStruct struct {
	Test string
}

func sendRequest(t *testing.T, router *evoli.Router, method string, path string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	return rr
}
