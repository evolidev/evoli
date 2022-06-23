package evoli

import (
	"github.com/evolidev/evoli/framework/console"
	"github.com/evolidev/evoli/framework/router"
	"github.com/evolidev/evoli/framework/use"
	"gorm.io/gorm"
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

func (a *Application) AddMigration(migrate func(db *gorm.DB)) {
	migrate(use.DB())
}

func (a *Application) Start() {
	cli := console.New()

	cli.AddCommand("routes", "List all registered routes", a.Serve)
	cli.AddCommand("make:routes", "List all registered routes", a.Serve)
	cli.AddCommand("make:controller", "List all registered routes", a.Serve)

	cli.Run()
}

func (a *Application) Serve(command *console.ParsedCommand) {
	log.Fatal(http.ListenAndServe(":8081", a.handler))
}

type MakeController struct {
	console.Command
}
