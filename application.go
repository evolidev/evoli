package evoli

import (
	"embed"
	"github.com/evolidev/evoli/framework/command"
	"github.com/evolidev/evoli/framework/component"
	"github.com/evolidev/evoli/framework/console"
	"github.com/evolidev/evoli/framework/logging"
	"github.com/evolidev/evoli/framework/middleware"
	"github.com/evolidev/evoli/framework/router"
	"github.com/evolidev/evoli/framework/use"
	"github.com/evolidev/evoli/framework/view"
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

	setupViewEngine()
	component.RegisterRoutes(handler)

	rootPath := use.BasePath()
	fsLogger := logging.NewLogger(&logging.Config{Name: "fs", PrefixColor: 32})
	fsLogger.Log("Setting root path to: " + rootPath)

	return &Application{
		handler: handler.AddMiddleware(middleware.NewLoggingMiddleware()),
		logger: logging.NewLogger(&logging.Config{
			Name:        "app",
			PrefixColor: 120,
		}),
	}
}

func setupViewEngine() {
	viewEngine := view.NewEngine()
	component.SetupViewEngine(viewEngine)

	use.AddFacade("viewEngine", viewEngine)
}

func (a *Application) AddRoutes(prefix string, routes func(router *router.Router)) {
	a.handler.Prefix(prefix).Group(routes)
}

func (a *Application) AddMigration(migrate func(db *gorm.DB)) {
	migrate(use.DB())
	l := logging.NewLogger(&logging.Config{Name: "db", PrefixColor: 50})
	l.Log("Models migrated successfully")
}

func (a *Application) Start() {
	cli := console.New()

	cli.AddCommand("serve {--port=8081}", "Serve the application", a.Serve)
	cli.Add(command.About())

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
