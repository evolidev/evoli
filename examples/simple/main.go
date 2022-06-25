package main

import (
	"embed"
	"github.com/evolidev/evoli"
	"github.com/evolidev/evoli/examples/simple/components"
	"github.com/evolidev/evoli/examples/simple/database"
	"github.com/evolidev/evoli/examples/simple/routes"
	"github.com/evolidev/evoli/framework/component"
	"github.com/evolidev/evoli/framework/logging"
)

var logger *logging.Logger
var app *evoli.Application

//go:embed resources public
var content embed.FS

func main() {
	logger = logging.NewLogger(&logging.Config{Name: "simple application", PrefixColor: 73})

	component.Register(components.Login{})
	component.Register(components.Input{})

	app = evoli.NewApplication()
	app.SetFS(content)
	app.AddRoutes("/", routes.Web)
	app.AddRoutes("/api", routes.Api)
	app.AddRoutes("/assets", routes.Files)
	app.AddMigration(database.Migrate)

	app.Start()
}

func consoleTest() {
	logger.Debug("Console test..D")
}

/**
 * HelloWorldWithPath component
 */

type HelloWorldWithPath struct {
}

func (h *HelloWorldWithPath) GetFilePath() string {
	return "hello-worlds"
}

func (h *HelloWorldWithPath) Test() string {
	logger.Debug("It is working!")
	return "super."
}

func helloWorldComponentTest() {
	component.Register(HelloWorldWithPath{})

	json := `{"Name":"Foo"}`

	hello := component.NewByNameWithData("HelloWorldWithPath", json)

	response := hello.Call("Test", nil)

	logger.Success(response)
}
