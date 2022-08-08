package main

import (
	"embed"
	"github.com/evolidev/evoli"
	"github.com/evolidev/evoli/examples/simple/components"
	"github.com/evolidev/evoli/examples/simple/database"
	"github.com/evolidev/evoli/examples/simple/routes"
	"github.com/evolidev/evoli/framework/component"
)

var app *evoli.Application

//go:embed resources public configs
var content embed.FS

func main() {

	component.Register(components.Login{})
	component.Register(components.Input{})

	app = evoli.NewApplication()
	app.AddEmbedFS(content)
	app.AddRoutes("/", routes.Web)
	app.AddRoutes("/api", routes.Api)
	app.AddRoutes("/assets", routes.Folders)
	app.AddRoutes("/", routes.Files)
	app.AddMigration(database.Migrate)

	app.Start()
}
