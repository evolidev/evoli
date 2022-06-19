package evoli

import (
	"github.com/evolidev/evoli/framework/router"
	"log"
	"net/http"
)

type App struct {
	handler *router.Router
}

func NewApplication() *App {
	return &App{
		handler: router.NewRouter(),
	}
}

func (a *App) AddRoutes(prefix string, routes func(router *router.Router)) {
	a.handler.Prefix(prefix).Group(routes)
}

func (a *App) Start() {
	log.Fatal(http.ListenAndServe(":8081", a.handler))
}

func Start() {
	//console.Commands()
	//watch()

}
