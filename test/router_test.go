package test

import (
	"encoding/json"
	"github.com/evolidev/evoli/framework/middleware"
	"github.com/evolidev/evoli/framework/response"
	evoli "github.com/evolidev/evoli/framework/router"
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

		path = "/response/views"
		router.Get(path, handlerViewResponse)

		rr = sendRequest(t, router, http.MethodGet, path)

		assert.Exactly(t, "<div>Hello test</div>", rr.Body.String())
	})

	t.Run("Basic route should have access to controller properties", func(t *testing.T) {
		path := "/controller"
		router.Get(path, MyController.TestAction)

		rr := sendRequest(t, router, http.MethodGet, path)

		assert.Exactly(t, path, rr.Body.String())
	})

	t.Run("Basic route should get parameter injected", func(t *testing.T) {
		path := "/controller/:param"
		router.Get(path, MyController.TestActionWithParam)

		rr := sendRequest(t, router, http.MethodGet, "/controller/test")

		assert.Exactly(t, "test", rr.Body.String())
	})

	t.Run("Basic route should get parameter injected", func(t *testing.T) {
		path := "/controller/:param"
		router.Get(path, MyController.TestActionWithParamAndRequest)

		rr := sendRequest(t, router, http.MethodGet, "/controller/test")

		assert.Exactly(t, "/controller/test/test", rr.Body.String())
	})

	t.Run("Basic route should get parameter injected in any order", func(t *testing.T) {
		path := "/controller/:param"
		router.Get(path, MyController.TestActionWithParamAndRequestOrdered)

		rr := sendRequest(t, router, http.MethodGet, "/controller/test")

		assert.Exactly(t, "/controller/test/test", rr.Body.String())
	})

	t.Run("Basic callback route should get parameter injected", func(t *testing.T) {
		path := "/inject/into/callback"
		router.Get(path, func(request *http.Request) string {
			return request.URL.Path
		})

		rr := sendRequest(t, router, http.MethodGet, "/inject/into/callback")

		assert.Exactly(t, "/inject/into/callback", rr.Body.String())
	})
}

func TestPrefix(t *testing.T) {
	t.Run("Prefix should prefix all sub routes", func(t *testing.T) {
		router := evoli.NewRouter()
		router.Prefix("/prefix").Group(func(router *evoli.Router) {
			router.Get("/test", func() string { return "prefix-test" })
		})

		rr := sendRequest(t, router, http.MethodGet, "/prefix/test")

		assert.Exactly(t, "prefix-test", rr.Body.String())
	})

	t.Run("Prefix should handle sub prefix routes too", func(t *testing.T) {
		router := evoli.NewRouter()
		router.Prefix("/prefix").Group(func(router *evoli.Router) {
			router.Prefix("/sub-prefix").Group(func(router *evoli.Router) {
				router.Get("/test", func() string { return "sub-prefix-test" })
			})
		})

		rr := sendRequest(t, router, http.MethodGet, "/prefix/sub-prefix/test")

		assert.Exactly(t, "sub-prefix-test", rr.Body.String())
	})

	t.Run("Prefix should not add prefix to outside handler", func(t *testing.T) {
		router := evoli.NewRouter()
		router.Prefix("/prefix").Group(func(router *evoli.Router) {
			router.Get("/test", func() string { return "sub-prefix-test" })
		})

		router.Get("/no-prefix", func() string { return "no-prefix-test" })

		rr := sendRequest(t, router, http.MethodGet, "/no-prefix")

		assert.Exactly(t, "no-prefix-test", rr.Body.String())
	})
}

func TestMiddleware(t *testing.T) {
	t.Run("Middleware should accept a handler func", func(t *testing.T) {
		router := evoli.NewRouter()

		var mid middleware.MiddlewareFunc
		mid = func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.Header.Set("my-test-header", "test")

				next.ServeHTTP(w, r)
			})
		}

		router.Middleware(mid).Group(func(router *evoli.Router) {
			router.Get("/middleware/test", func(request *http.Request) string { return request.Header.Get("my-test-header") })
		})

		rr := sendRequest(t, router, http.MethodGet, "/middleware/test")

		assert.Exactly(t, "test", rr.Body.String())
	})

	t.Run("Middleware should accept a handler", func(t *testing.T) {
		router := evoli.NewRouter()

		router.Middleware(MyMiddleware{testHeader: "test"}).Group(func(router *evoli.Router) {
			router.Get("/middleware/test", func(request *http.Request) string { return request.Header.Get("my-test-header") })
		})

		rr := sendRequest(t, router, http.MethodGet, "/middleware/test")

		assert.Exactly(t, "test", rr.Body.String())
	})
}

type MyMiddleware struct {
	testHeader string
}

func (m MyMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		request.Header.Set("my-test-header", m.testHeader)

		next.ServeHTTP(writer, request)
	})
}

type MyController struct {
	Request http.Request
}

func (m MyController) TestActionWithParam(test string) string {
	return test
}

func (m MyController) TestActionWithParamAndRequest(request *http.Request, test string) string {
	return request.URL.Path + "/" + test
}

func (m MyController) TestActionWithParamAndRequestOrdered(test string, request *http.Request) string {
	return request.URL.Path + "/" + test
}

func (m MyController) TestAction() string {
	return m.Request.URL.Path
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

func handlerViewResponse() interface{} {
	return response.View("templates.test")
}

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
