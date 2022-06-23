package evoli

import (
	"embed"
	"github.com/evolidev/evoli/framework/router"
	"github.com/evolidev/evoli/framework/use"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Application struct {
	handler *router.Router
	fs      embed.FS
}

func NewApplication() *Application {
	return &Application{
		handler: router.NewRouter(),
	}
}

func (a *Application) AddRoutes(prefix string, routes func(router *router.Router)) {
	a.handler.Prefix(prefix).Group(routes)
}

func (a *Application) AddMigration(migrate func(db *gorm.DB)) {
	migrate(use.DB())
}

func (a *Application) Start() {
	log.Fatal(http.ListenAndServe(":8081", a.handler))
}

func (a *Application) SetFS(fs embed.FS) {
	a.fs = fs
	a.handler.Fs = a.fs
}

func Start() {
	//console.Commands()
	//watch()

}
