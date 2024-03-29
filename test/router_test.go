package test

import (
	"bytes"
	"encoding/json"
	"github.com/evolidev/evoli/framework/config"
	"github.com/evolidev/evoli/framework/middleware"
	"github.com/evolidev/evoli/framework/response"
	evoli "github.com/evolidev/evoli/framework/router"
	"github.com/evolidev/evoli/framework/use"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
)

func TestBasic(t *testing.T) {
	t.Parallel()
	router := evoli.NewRouter()
	use.Embed(tmp)

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
		path := "/controller/only/:param"
		router.Get(path, MyController.TestActionWithParam)

		rr := sendRequest(t, router, http.MethodGet, "/controller/only/test")

		assert.Exactly(t, "test", rr.Body.String())
	})

	t.Run("Basic route should get parameter injected", func(t *testing.T) {
		path := "/controller/injected/:param"
		router.Get(path, MyController.TestActionWithParamAndRequest)

		rr := sendRequest(t, router, http.MethodGet, "/controller/injected/test")

		assert.Exactly(t, "/controller/injected/test/test", rr.Body.String())
	})

	t.Run("Basic route should get parameter injected in any order", func(t *testing.T) {
		path := "/controller/injected-any-order/:param"
		router.Get(path, MyController.TestActionWithParamAndRequestOrdered)

		rr := sendRequest(t, router, http.MethodGet, "/controller/injected-any-order/test")

		assert.Exactly(t, "/controller/injected-any-order/test/test", rr.Body.String())
	})

	t.Run("Basic callback route should get parameter injected", func(t *testing.T) {
		path := "/inject/into/callback"
		router.Get(path, func(request *http.Request) string {
			return request.URL.Path
		})

		rr := sendRequest(t, router, http.MethodGet, "/inject/into/callback")

		assert.Exactly(t, "/inject/into/callback", rr.Body.String())
	})

	t.Run("Basic callback routes should get route params injected", func(t *testing.T) {
		path := "/inject/route/param/:test"
		router.Get(path, func(params httprouter.Params) string {
			return params.ByName("test")
		})

		rr := sendRequest(t, router, http.MethodGet, "/inject/route/param/awesome")

		assert.Exactly(t, "awesome", rr.Body.String())
	})

	t.Run("Basic callback routes should get query params injected", func(t *testing.T) {
		path := "/inject/route/query"
		router.Get(path, func(params url.Values) string {
			return params.Get("test")
		})

		rr := sendRequest(t, router, http.MethodGet, "/inject/route/query?test=awesome")

		assert.Exactly(t, "awesome", rr.Body.String())
	})

	t.Run("Basic callback with no response should return no data header", func(t *testing.T) {
		path := "/basic/no-return"
		router.Get(path, func() {})

		rr := sendRequest(t, router, http.MethodGet, path)

		assert.Exactly(t, http.StatusNoContent, rr.Code)
	})
}

func TestBasicPost(t *testing.T) {
	router := evoli.NewRouter()
	t.Run("Basic post should map form data to struct", func(t *testing.T) {
		path := "/basic/post/form"
		router.Post(path, func(inputStruct InputStruct) string {
			return inputStruct.Test
		})

		f := url.Values{}
		f.Add("test", "my-form")
		rr := sendRequestWithForm(t, router, http.MethodGet, path, f)

		assert.Exactly(t, "my-form", rr.Body.String())
	})

	t.Run("Basic post with url params should map form data to struct", func(t *testing.T) {
		path := "/basic/post/with-param/:param"
		router.Post(path, func(param string, inputStruct InputStruct) string {
			return param + " " + inputStruct.Test
		})

		f := url.Values{}
		f.Add("test", "my-form")
		rr := sendRequestWithForm(t, router, http.MethodGet, "/basic/post/with-param/test", f)

		assert.Exactly(t, "test my-form", rr.Body.String())
	})
}

func TestPrefix(t *testing.T) {
	use.Embed(tmp)

	t.Parallel()
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
	use.Embed(tmp)

	t.Parallel()
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

func TestRedirectResponse(t *testing.T) {
	use.Embed(tmp)

	//t.Parallel()
	t.Run("A redirect response should redirect to desired route", func(t *testing.T) {
		router := evoli.NewRouter()

		router.Get("/redirect/response", func() *response.RedirectResponse {
			return response.Redirect("/redirect/to")
		})

		rr := sendRequest(t, router, http.MethodGet, "/redirect/response")

		assert.Equal(t, "/redirect/to", rr.Header().Get("Location"))
		assert.Exactly(t, http.StatusSeeOther, rr.Code)
	})
}

func TestRequest(t *testing.T) {
	use.Embed(tmp)

	t.Parallel()
	t.Run("Request struct should return route param by name", func(t *testing.T) {
		r := evoli.NewRouter()

		r.Get("/request/route/param/:test", func(request *evoli.Request) string {
			return request.Get("test").(string)
		})

		rr := sendRequest(t, r, http.MethodGet, "/request/route/param/test")

		assert.Equal(t, "test", rr.Body.String())
	})

	t.Run("Request struct should return query param by name", func(t *testing.T) {
		r := evoli.NewRouter()

		r.Get("/request/route/param/query", func(request *evoli.Request) string {
			return request.Get("test").(string)
		})

		rr := sendRequest(t, r, http.MethodGet, "/request/route/param/query?test=awesome")

		assert.Equal(t, "awesome", rr.Body.String())
	})
	t.Run("Request struct should return form param by name", func(t *testing.T) {
		r := evoli.NewRouter()
		f := url.Values{}
		f.Add("form", "my-form")
		r.Post("/request/form/param", func(request *evoli.Request) string {
			return request.Get("form").(string)
		})

		rr := sendRequestWithForm(t, r, http.MethodGet, "/request/form/param", f)

		assert.Equal(t, "my-form", rr.Body.String())
	})
	t.Run("Request struct should return body param by name", func(t *testing.T) {
		r := evoli.NewRouter()
		f := make(map[string]string)
		f["body"] = "myBody"
		r.Get("/request/body/param", func(request *evoli.Request) string {
			return request.Get("body").(string)
		})

		rr := sendRequestWithData(t, r, http.MethodGet, "/request/body/param", f)

		assert.Equal(t, "myBody", rr.Body.String())
	})
}

func TestStatic(t *testing.T) {
	t.Run("static should serve static files in folder", func(t *testing.T) {
		config.SetDirectory("./")
		cnf := use.Config("storage")
		abs, _ := filepath.Abs("./")
		cnf.Set("local.base_path", abs)

		r := evoli.NewRouter()

		r.Static("/resources/local", "resources/files")

		rr := sendRequest(t, r, http.MethodGet, "/resources/local/test.txt")

		assert.Equal(t, "test", rr.Body.String())
	})

	t.Run("static should serve static files in folder in right storage", func(t *testing.T) {
		use.Embed(tmp)
		config.SetDirectory("./")
		cnf := use.Config("storage")
		abs, _ := filepath.Abs("./")
		cnf.Set("local.base_path", abs)

		r := evoli.NewRouter()

		r.Static("/resources/static/local", "local:static")

		rr := sendRequest(t, r, http.MethodGet, "/resources/static/local/test.txt")

		assert.Equal(t, "test", rr.Body.String())
	})
}

type InputStruct struct {
	Test string
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
	HttpRequest http.Request
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
	return m.HttpRequest.URL.Path
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
	req := httptest.NewRequest(method, path, nil)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	return rr
}

func sendRequestWithForm(t *testing.T, router *evoli.Router, method string, path string, form url.Values) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	req.Form = form
	req.Method = http.MethodPost

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	return rr
}

func sendRequestWithData(t *testing.T, router *evoli.Router, method string, path string, data any) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(use.JsonEncode(data)))

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	return rr
}
