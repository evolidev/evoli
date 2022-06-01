package evoli

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Router struct {
	router httprouter.Router
}

func (r *Router) Get(path string, f func() string) {
	r.router.GET(path, func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		writer.Write([]byte(f()))
	})
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func NewRouter() *Router {
	return &Router{}
}
