package evoli

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Router struct {
	router httprouter.Router
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *Router) Get(path string, f func() string) {
	r.handle(http.MethodGet, path, f)
}

func (r *Router) Post(path string, f func() string) {
	r.handle(http.MethodPost, path, f)
}

func (r *Router) Put(path string, f func() string) {
	r.handle(http.MethodPut, path, f)
}

func (r *Router) Patch(path string, f func() string) {
	r.handle(http.MethodPatch, path, f)
}

func (r *Router) Delete(path string, f func() string) {
	r.handle(http.MethodDelete, path, f)
}

func (r *Router) Options(path string, f func() string) {
	r.handle(http.MethodOptions, path, f)
}

func (r *Router) Head(path string, f func() string) {
	r.handle(http.MethodHead, path, f)
}

func (r *Router) handle(method string, path string, f func() string) {
	r.router.Handle(method, path, func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		writer.Write([]byte(f()))
	})
}

func (r *Router) Match(path string, f func() string, httpMethods ...string) {
	for _, httpMehtod := range httpMethods {
		r.handle(httpMehtod, path, f)
	}
}

func NewRouter() *Router {
	r := &Router{}
	r.router.HandleMethodNotAllowed = true

	return r
}
