package command

import (
	"github.com/evolidev/evoli/framework/console"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"strings"
	"text/template"
)

func Route() *console.Command {
	return &console.Command{
		Definition:  "routes {name}",
		Description: "Add a new routes file",
		Execution:   routeRun,
	}
}

func routeRun(cmd *console.ParsedCommand) {
	os.Mkdir("routes", os.ModePerm)
	os.Mkdir("resources", os.ModePerm)
	os.Mkdir("resources/views", os.ModePerm)

	name := cmd.GetArgument("name").String()

	f, _ := os.Create("routes/" + name + ".go")
	v, _ := os.Create("resources/views/my-endpoint.html")
	defer f.Close()
	defer v.Close()

	data := RouteData{
		ParamName: strings.ToLower(name),
		Name:      cases.Title(language.English, cases.Compact).String(name),
	}

	routeTemplate.Execute(f, data)
	viewTemplate.Execute(v, data)

	generate()
}

type RouteData struct {
	Name      string
	ParamName string
}

var viewTemplate = template.Must(template.New("").Parse(`
<div>Hiho</div>
`))

var routeTemplate = template.Must(template.New("").Parse(`
package routes

import (
	"github.com/evolidev/evoli/framework/response"
	"github.com/evolidev/evoli/framework/router"
)

func {{.Name}}({{.ParamName}} *router.Router) {
	{{.ParamName}}.Get("/", func() string { return "hello from {{.Name}}" })
	{{.ParamName}}.Get("/my-endpoint", func() *response.ViewResponse {
		return response.View("my-endpoint")
	})
}
`))
