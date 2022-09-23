package command

import (
	"fmt"
	"github.com/evolidev/evoli/framework/console"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"text/template"
)

func Model() *console.Command {
	return &console.Command{
		Definition:  "models {name}",
		Description: "Creates a new model",
		Execution:   modelRun,
	}
}

type ModelData struct {
	Name    string
	Props   map[string]string
	Imports []string
}

func modelRun(cmd *console.ParsedCommand) {
	fmt.Println("hiho")
	name := cmd.GetArgument("name").String()
	props := cmd.GetArgument("props")

	fmt.Println(props)

	os.Mkdir("models", os.ModePerm)
	f, _ := os.Create("models/" + name + ".go")

	data := ModelData{
		Name:    cases.Title(language.English, cases.Compact).String(name),
		Props:   make(map[string]string),
		Imports: make([]string, 0),
	}

	if name == "user" {
		data.Props["Email"] = "string"
		data.Props["Password"] = "string"
	}

	if name == "session" {
		data.Props["User"] = "User"
		data.Props["IpAddress"] = "string"
	}

	modelTemplate.Execute(f, data)

}

var modelTemplate = template.Must(template.New("").Parse(`
package models

import (
	"gorm.io/gorm"
)

type {{.Name}} struct {
	gorm.Model
	{{- range $prop, $type := .Props}}
	{{$prop}} {{$type}}
	{{- end}}
}
`))
