package evoli

import (
	"github.com/evolidev/evoli/framework/logging"
	"github.com/evolidev/evoli/framework/response"
	"github.com/evolidev/evoli/framework/use"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

type Router struct {
	router *httprouter.Router
	prefix string
	logger *logging.Logger
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
	if r.prefix != "/" {
		if path == "/" {
			path = ""
		}
		path = r.prefix + path
	}

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
	router := &Router{router: httprouter.New(), prefix: "/"}

	router.router.RedirectTrailingSlash = false
	router.router.RedirectFixedPath = false

	return router
}

func NewRouteSwitch() *RouteSwitch {
	return &RouteSwitch{
		routes: use.NewCollection[string, *Router](),
	}
}

type RouteSwitch struct {
	routes *use.Collection[string, *Router]
	logger *logging.Logger
}

func (rs *RouteSwitch) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	base := "/"
	router := rs.routes.Get(base)

	path := request.URL.Path

	split := strings.Split(path, "/")

	if len(split) > 1 {
		base = "/" + split[1]
	}

	if rs.routes.Has(base) {
		router = rs.routes.Get(base)
	}

	router.ServeHTTP(writer, request)
}

func (rs *RouteSwitch) Get(prefix string) *Router {
	return rs.routes.Get(prefix)
}

func (rs *RouteSwitch) Add(prefix string, routes func(router *Router)) {
	router := NewRouter()
	router.prefix = prefix
	routes(router)

	rs.routes.Add(prefix, router)
}
