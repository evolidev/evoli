package router

import (
	"embed"
	"fmt"
	"github.com/evolidev/evoli/framework/logging"
	"github.com/evolidev/evoli/framework/middleware"
	"github.com/evolidev/evoli/framework/response"
	"github.com/evolidev/evoli/framework/use"
	"github.com/julienschmidt/httprouter"
	"io/fs"
	"net/http"
	"strings"
)

type Router struct {
	router      *httprouter.Router
	prefix      string
	logger      *logging.Logger
	middlewares []middleware.Middleware
	Fs          embed.FS
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
	var next http.Handler
	next = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		req := NewRequest(request)
		tmp := use.Magic(handler)

		myResponse := response.NewResponse(
			tmp.WithParams(req.Params()).Fill().WithParams(req.Params().Get("Route")).Call(),
		)

		myResponse.Headers().Iterate(func(key string, value string) {
			writer.Header().Add(key, value)
		})

		if redirect, ok := myResponse.(*response.RedirectResponse); ok {
			http.Redirect(writer, request, redirect.To, myResponse.Code())

			return
		}

		responseCode := myResponse.Code()
		if responseCode == 0 {
			responseCode = http.StatusOK
		}
		writer.WriteHeader(responseCode)

		_, werr := writer.Write(myResponse.AsBytes())
		if werr != nil {
			// todo log to our logger and redirect to 500?
			fmt.Println(werr)
		}
	})

	for _, m := range r.middlewares {
		next = m.Middleware(next)
	}

	r.router.Handler(method, r.pathWithPrefix(path), next)
}

func (r *Router) ServeFiles(path string, fs http.FileSystem) {
	r.router.ServeFiles(r.pathWithPrefix(path)+"/*filepath", fs)
}

func (r *Router) pathWithPrefix(path string) string {
	if r.prefix != "/" {
		if path == "/" {
			path = ""
		}
		path = r.prefix + path
	}

	return path
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

func (r *Router) Prefix(prefix string) *Group {
	group := NewGroup(r)
	group.router.prefix = r.prefix + prefix

	group.router.prefix = strings.Replace(group.router.prefix, "//", "/", 1)

	return group
}

func (r *Router) AddMiddleware(middleware middleware.Middleware) *Router {
	r.middlewares = append(r.middlewares, middleware)

	return r
}

func (r *Router) Middleware(middlewares ...middleware.Middleware) *Group {
	group := NewGroup(r)
	group.router.middlewares = append(group.router.middlewares, middlewares...)

	//group.router.prefix = r.prefix + prefix
	//
	//group.router.prefix = strings.Replace(group.router.prefix, "//", "/", 1)

	return group
}

func (r *Router) Static(path string, rootDir string) {
	_, err := fs.ReadDir(r.Fs, rootDir)
	var servingFS http.FileSystem

	if err != nil {
		servingFS = http.FileSystem(http.Dir(rootDir))
	} else {
		sub, _ := fs.Sub(r.Fs, rootDir)
		servingFS = http.FS(sub)
	}

	r.ServeFiles(path, servingFS)
}

func (r *Router) File(path string, file string) {
	r.router.Handler(http.MethodGet, r.pathWithPrefix(path), http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, file)
	}))
}

func NewRouter() *Router {
	router := &Router{router: httprouter.New(), prefix: "/"}

	router.router.RedirectTrailingSlash = false
	router.router.RedirectFixedPath = false

	router.middlewares = make([]middleware.Middleware, 0)

	return router
}

type Group struct {
	router *Router
}

func (g *Group) Group(routes func(*Router)) {
	routes(g.router)
}

func NewGroup(router *Router) *Group {
	return &Group{router: &Router{router: router.router, middlewares: router.middlewares, Fs: router.Fs}}
}
