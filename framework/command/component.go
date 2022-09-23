package command

import (
	"github.com/evolidev/evoli/framework/console"
	"github.com/evolidev/evoli/framework/use"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"strings"
	"text/template"
)

func Component() *console.Command {
	return &console.Command{
		Definition:  "component {name}",
		Description: "Add a new component",
		Execution:   componentRun,
	}
}

func componentRun(cmd *console.ParsedCommand) {
	os.Mkdir("components", os.ModePerm)
	os.Mkdir("resources", os.ModePerm)
	os.Mkdir("resources/views/components", os.ModePerm)

	name := cmd.GetArgument("name").String()

	data := ComponentData{
		ParamName: use.String(name).Kebab().Get(),
		Name:      cases.Title(language.English, cases.Compact).String(name),
		Short:     strings.ToLower(name)[0:1],
	}

	f, _ := os.Create("components/" + data.Name + ".go")
	v, _ := os.Create("resources/views/components/" + data.ParamName + ".html")
	defer f.Close()
	defer v.Close()

	componentTemplate.Execute(f, data)
	componentViewTemplate.Execute(v, data)

	generate()
}

type ComponentData struct {
	Name      string
	ParamName string
	Short     string
}

var componentTemplate = template.Must(template.New("").Parse(`
package components

type {{.Name}} struct {
	Name string
}

func ({{.Short}} *{{.Name}}) Mount() {
	//todo do stuff here on mount
}

func ({{.Short}} *{{.Name}}) Update(p any) {
	{{.Short}}.Name = p.(string)
}
`))

var componentViewTemplate = template.Must(template.New("").Parse(`
<div @scope>
${ .Name }
</div>
`))
