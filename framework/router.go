package evoli

import (
	"github.com/evolidev/evoli/framework/response"
	"github.com/evolidev/evoli/framework/use"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Router struct {
	router *httprouter.Router
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *Router) Get(path string, handler interface{}) {
	r.handle(http.MethodGet, path, handler)
}

func (r *Router) Post(path string, handler interface{}) {
	r.handle(http.MethodPost, path, handler)
}

func (r *Router) Put(path string, handler interface{}) {
	r.handle(http.MethodPut, path, handler)
}

func (r *Router) Patch(path string, handler interface{}) {
	r.handle(http.MethodPatch, path, handler)
}

func (r *Router) Delete(path string, handler interface{}) {
	r.handle(http.MethodDelete, path, handler)
}

func (r *Router) Options(path string, handler interface{}) {
	r.handle(http.MethodOptions, path, handler)
}

func (r *Router) Head(path string, handler interface{}) {
	r.handle(http.MethodHead, path, handler)
}

func (r *Router) Connect(path string, handler interface{}) {
	r.handle(http.MethodConnect, path, handler)
}

func (r *Router) Trace(path string, handler interface{}) {
	r.handle(http.MethodTrace, path, handler)
}

func (r *Router) handle(method string, path string, handler interface{}) {
	r.router.Handle(method, path, func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		tmp := use.Magic(handler)

		myParams := make(map[string]interface{}, 0)
		myParams["Request"] = request

		response := response.NewResponse(tmp.WithParams(myParams).Fill().WithParams(params).Call())

		response.Headers().Iterate(func(key string, value string) {
			writer.Header().Add(key, value)
		})

		writer.Header().Add("Content-Type", "charset=utf-8")

		writer.Write(response.AsBytes())
	})
}

func (r *Router) Match(path string, handler interface{}, httpMethods ...string) {
	for _, httpMethod := range httpMethods {
		r.handle(httpMethod, path, handler)
	}
}

func (r *Router) Any(path string, handler interface{}) {
	for _, fn := range r.MethodTable() {
		fn(path, handler)
	}
}

func (r *Router) MethodTable() map[string]func(path string, handler interface{}) {
	methodTable := make(map[string]func(path string, handler interface{}))
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
