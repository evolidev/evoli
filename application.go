package evoli

import (
	"embed"
	"github.com/evolidev/evoli/framework/console"
	"github.com/evolidev/evoli/framework/logging"
	"github.com/evolidev/evoli/framework/middleware"
	"github.com/evolidev/evoli/framework/router"
	"github.com/evolidev/evoli/framework/use"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Application struct {
	handler *router.Router
	logger  *logging.Logger
	fs      embed.FS
}

func NewApplication() *Application {
	handler := router.NewRouter()

	return &Application{
		handler: handler.AddMiddleware(middleware.NewLoggingMiddleware()),
		logger: logging.NewLogger(&logging.Config{
			Name:        "app",
			PrefixColor: 120,
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

	cli.Run()
}

func (a *Application) SetFS(fs embed.FS) {
	a.fs = fs
	a.handler.Fs = a.fs
}

func (a *Application) Serve(command *console.ParsedCommand) {
	port := command.GetOptionWithDefault("port", 8081).(string)

	a.logger.Log("Serving application on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, a.handler))
}

type MakeController struct {
	console.Command
}
