package evoli

import (
	"github.com/evolidev/evoli/framework/router"
	"log"
	"net/http"
)

type Application struct {
	handler *router.Router
}

func NewApplication() *Application {
	return &Application{
		handler: router.NewRouter(),
	}
}

func (a *Application) AddRoutes(prefix string, routes func(router *router.Router)) {
	a.handler.Prefix(prefix).Group(routes)
}

func (a *Application) Start() {
	log.Fatal(http.ListenAndServe(":8081", a.handler))
}

func Start() {
	//console.Commands()
	//watch()

}
