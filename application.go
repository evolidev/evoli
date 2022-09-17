package evoli

import (
	"embed"
	"fmt"
	"github.com/evolidev/evoli/framework/command"
	"github.com/evolidev/evoli/framework/component"
	"github.com/evolidev/evoli/framework/console"
	"github.com/evolidev/evoli/framework/console/reload"
	"github.com/evolidev/evoli/framework/filesystem"
	"github.com/evolidev/evoli/framework/logging"
	"github.com/evolidev/evoli/framework/middleware"
	"github.com/evolidev/evoli/framework/router"
	"github.com/evolidev/evoli/framework/use"
	"github.com/evolidev/evoli/framework/view"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

//go:embed1 resources
//var evoliFs embed.FS

type Application struct {
	handler *router.Router
	logger  *logging.Logger
	Cli     *console.Console
}

func NewApplication() *Application {
	app := &Application{}
	app.Init()

	return app
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
	use.Migration().Add(migrate)
}

func (a *Application) Start() {
	a.listenForSignal()

	a.Cli.AddCommand("serve {--port=8081}", "Serve the application", a.Serve)
	a.Cli.AddCommand("watch {--port=8081}", "Serve and watch the application", a.Watch)
	a.Cli.Add(command.About())
	a.Cli.Add(command.Migrate())
	a.Cli.Add(command.Generate())
	a.Cli.Add(command.Init())
	a.Cli.Add(command.Route())
	a.Cli.Add(command.Component())
	a.Cli.Add(command.Model())

	a.Cli.Run()
}

func (a *Application) Init() {
	handler := router.NewRouter()
	handler = handler.AddMiddleware(middleware.NewLoggingMiddleware())

	//oldFs := handler.Fs
	//handler.Fs = evoliFs
	//handler.Static("/vendor/evoli/static", "resources")
	//handler.Fs = oldFs

	setupViewEngine()
	component.RegisterRoutes(handler)

	use.BasePath()

	logger := logging.NewLogger(&logging.Config{
		Name:        "app",
		PrefixColor: 32,
	})

	logging.SetAppLogger(logger)

	a.handler = handler
	a.logger = logger
	a.Cli = console.New()
}

func (a *Application) RegisterComponent(comp component.Component) {
	component.Register(comp)
}

func (a *Application) AddEmbedFS(fs embed.FS) {
	use.Embed(fs)
}

func (a *Application) Watch(command *console.ParsedCommand) {
	config := &reload.Configuration{
		AppRoot:            use.BasePath(),
		IncludedExtensions: []string{".go"},
		BuildPath:          "",
		BinaryName:         "main.go",
		Command:            "go run main.go serve",
		Debug:              false,
		ForcePolling:       false,
	}

	reload.RunBackground(config)
}

func (a *Application) Serve(command *console.ParsedCommand) {
	autoMigrateIfEnabled()

	port := command.GetOptionWithDefault("port", 8081).(string)

	filesystem.Write(use.StoragePath("tmp/serve.pid"), fmt.Sprintf("%d", os.Getpid()))
	defer filesystem.Delete(use.StoragePath("tmp/serve.pid"))

	a.logger.Log("Serving application on http://localhost:%s (PID: %d)", port, os.Getpid())
	a.logger.Fatal(http.ListenAndServe(":"+port, a.handler))
}

func autoMigrateIfEnabled() {
	if use.Config("db.auto_migrate").Value().(bool) {
		use.Migration().Migrate(use.DB())
	}
}

func (a *Application) listenForSignal() {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		s := <-sigChannel
		a.logger.Debug("received signal: %s", s)
		os.Exit(0)
	}()
}
