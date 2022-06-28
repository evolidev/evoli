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
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	use.BasePath()

	return &Application{
		handler: handler.AddMiddleware(middleware.NewLoggingMiddleware()),
		logger: logging.NewLogger(&logging.Config{
			Name:        "app",
			PrefixColor: 32,
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
	a.listenForSignal()

	cli := console.New()

	cli.AddCommand("serve {--port=8081}", "Serve the application", a.Serve)
	cli.AddCommand("watch {--port=8081}", "Serve and watch the application", a.Watch)
	cli.Add(command.About())

	cli.Run()
}

func (a *Application) SetFS(fs embed.FS) {
	a.fs = fs
	a.handler.Fs = a.fs
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
	port := command.GetOptionWithDefault("port", 8081).(string)

	filesystem.Write(use.StoragePath("tmp/serve.pid"), fmt.Sprintf("%d", os.Getpid()))
	defer filesystem.Delete(use.StoragePath("tmp/serve.pid"))

	a.logger.Log("Serving application on http://localhost:%s (PID: %d)", port, os.Getpid())
	log.Fatal(http.ListenAndServe(":"+port, a.handler))
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
