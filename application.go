package evoli

import (
	"github.com/evolidev/evoli/framework/console"
	"github.com/evolidev/evoli/framework/logging"
	"github.com/evolidev/evoli/framework/router"
	"github.com/evolidev/evoli/framework/use"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Application struct {
	handler *router.Router
	logger  *logging.Logger
}

func NewApplication() *Application {
	return &Application{
		handler: router.NewRouter(),
		logger: logging.NewLogger(&logging.Config{
			Name: "app",
		}),
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

	cli.AddCommand("serve {--port=8081}", "Serve the application", a.Serve)

	cli.AddCommand("route:list", "List all registered routes", a.Serve)
	cli.AddCommand("make:routes", "List all registered routes", a.Serve)
	cli.AddCommand("make:controller", "List all registered routes", a.Serve)

	cli.Run()
}

func (a *Application) Serve(command *console.ParsedCommand) {
	port := command.GetOption("port").(string)
	if port == "" {
		port = "8081"
	}

	a.logger.Log("Serving application on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, a.handler))
}

type MakeController struct {
	console.Command
}
