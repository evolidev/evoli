package evoli

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type Router struct {
	router *httprouter.Router
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *Router) Get(path string, f func() interface{}) {
	r.handle(http.MethodGet, path, f)
}

func (r *Router) Post(path string, f func() interface{}) {
	r.handle(http.MethodPost, path, f)
}

func (r *Router) Put(path string, f func() interface{}) {
	r.handle(http.MethodPut, path, f)
}

func (r *Router) Patch(path string, f func() interface{}) {
	r.handle(http.MethodPatch, path, f)
}

func (r *Router) Delete(path string, f func() interface{}) {
	r.handle(http.MethodDelete, path, f)
}

func (r *Router) Options(path string, f func() interface{}) {
	r.handle(http.MethodOptions, path, f)
}

func (r *Router) Head(path string, f func() interface{}) {
	r.handle(http.MethodHead, path, f)
}

func (r *Router) Connect(path string, f func() interface{}) {
	r.handle(http.MethodConnect, path, f)
}

func (r *Router) Trace(path string, f func() interface{}) {
	r.handle(http.MethodTrace, path, f)
}

func (r *Router) handle(method string, path string, f func() interface{}) {
	r.router.Handle(method, path, func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		result := f()
		var response []byte
		contentType := "text/plain"
		switch result.(type) {
		case int:
			response = []byte(strconv.Itoa(result.(int)))
		case string:
			response = []byte(result.(string))
		default:
			// todo handle error
			response, _ = json.Marshal(result)
			contentType = "application/json"
		}

		writer.Header().Add("Content-Type", contentType)
		writer.Header().Add("Content-Type", "charset=utf-8")

		writer.Write(response)
	})
}

func (r *Router) Match(path string, f func() interface{}, httpMethods ...string) {
	for _, httpMethod := range httpMethods {
		r.handle(httpMethod, path, f)
	}
}

func (r *Router) Any(path string, f func() interface{}) {
	for _, fn := range r.MethodTable() {
		fn(path, f)
	}
}

func (r *Router) MethodTable() map[string]func(path string, f func() interface{}) {
	methodTable := make(map[string]func(path string, f func() interface{}))
	methodTable[http.MethodGet] = r.Get
	methodTable[http.MethodPost] = r.Post
	methodTable[http.MethodPut] = r.Put
	methodTable[http.MethodPatch] = r.Patch
	methodTable[http.MethodDelete] = r.Delete
	methodTable[http.MethodHead] = r.Head
	methodTable[http.MethodOptions] = r.Options
	methodTable[http.MethodConnect] = r.Connect
	methodTable[http.MethodTrace] = r.Trace

	return methodTable
}

func NewRouter() *Router {
	return &Router{router: httprouter.New()}
}
