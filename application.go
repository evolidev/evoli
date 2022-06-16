package evoli

import (
	evoli "github.com/evolidev/evoli/framework"
	"log"
	"net/http"
)

type App struct {
	handler *evoli.RouteSwitch
}

func NewApplication() *App {
	return &App{
		handler: evoli.NewRouteSwitch(),
	}
}

func (a *App) AddRoutes(prefix string, routes func(router *evoli.Router)) {
	a.handler.Add(prefix, routes)
}

func (a *App) Start() {
	log.Fatal(http.ListenAndServe(":8081", a.handler))
}

func Start() {
	//console.Commands()
	//watch()

}
