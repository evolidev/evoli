package test

import (
	"encoding/json"
	evoli "github.com/evolidev/evoli/framework"
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

		rr = sendRequest(t, router, http.MethodGet, pathSlice)

		testJson, err = json.Marshal([]uint8{255, 255, 255})
		if err != nil {
			t.Fatal(err)
		}

		assert.Exactly(t, string(testJson), rr.Body.String())
	})
}

func handler() interface{} {
	return "hello-world"
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
