package command

import (
	"fmt"
	"github.com/evolidev/evoli/framework/console"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/exp/slices"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

func Generate() *console.Command {
	return &console.Command{
		Definition:  "generate",
		Description: "Generates app.go to register all stuff",
		Execution:   generateRun,
	}
}

type Register struct {
	Pkg        string
	Timestamp  time.Time
	Routes     map[string]string
	Imports    []string
	Components []string
	App        string
	Accessor   string
}

func generateRun(cmd *console.ParsedCommand) {
	f, _ := os.Create("app.go")
	defer f.Close()

	baseDir, _ := os.Getwd()
	pathSplit := strings.Split(baseDir, "/")

	moduleName := pathSplit[len(pathSplit)-1]

	reg := Register{
		Routes:     make(map[string]string),
		Timestamp:  time.Now(),
		Imports:    make([]string, 0),
		Components: make([]string, 0),
		Accessor:   "Application",
	}

	reg.Imports = append(reg.Imports, "github.com/evolidev/evoli")
	reg.Pkg = "main"

	dir, _ := os.ReadDir("routes")
	for _, d := range dir {
		i := parse("routes/" + d.Name())
		for _, tmp := range i.funcs {
			routePath := strings.Split(d.Name(), ".")[0]
			routePath = strings.ToLower(routePath)

			if routePath == "web" || routePath == "static" {
				routePath = ""
			}
			reg.Routes[i.pkg+"."+tmp] = "/" + routePath
			if !slices.Contains(reg.Imports, moduleName+"/"+i.pkg) {
				reg.Imports = append(reg.Imports, moduleName+"/"+i.pkg)
			}
		}
	}

	dir, _ = os.ReadDir("components")
	reg.Imports = append(reg.Imports, moduleName+"/components")
	for _, d := range dir {
		componentName := strings.Split(d.Name(), ".")[0]
		componentName = strings.ToLower(componentName)
		componentName = cases.Title(language.English, cases.Compact).String(componentName)

		reg.Components = append(reg.Components, "components."+componentName+"{}")
	}

	i := parse("main.go")

	if i.accessor != "" {
		reg.Accessor = i.accessor
	}
	reg.App = i.structs[0]

	packageTemplate.Execute(f, reg)
}

type info struct {
	pkg      string
	funcs    []string
	structs  []string
	accessor string
}

func parse(fileName string) info {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	i := info{funcs: make([]string, 0), structs: make([]string, 0)}

	i.pkg = node.Name.Name

	for _, f := range node.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}

		i.funcs = append(i.funcs, fn.Name.Name)
	}

	for _, n := range node.Decls {
		switch n.(type) {

		case *ast.GenDecl:
			genDecl := n.(*ast.GenDecl)
			for _, spec := range genDecl.Specs {
				switch spec.(type) {
				case *ast.TypeSpec:
					typeSpec := spec.(*ast.TypeSpec)

					i.structs = append(i.structs, typeSpec.Name.Name)

					switch typeSpec.Type.(type) {
					case *ast.StructType:
						structType := typeSpec.Type.(*ast.StructType)
						for _, field := range structType.Fields.List {
							switch field.Type.(type) {
							case *ast.Ident:
								k := field.Type.(*ast.Ident)
								fieldType := k.Name

								for _, name := range field.Names {
									fmt.Printf("\tField: name=%s type=%s\n", name.Name, fieldType)
								}
							case *ast.StarExpr:
								k := field.Type.(*ast.StarExpr)
								fieldType := k.X
								//fmt.Println(fieldType == "&{evoli Application}")
								s := fmt.Sprintf("%s", fieldType)
								if s == "&{evoli Application}" {
									if len(field.Names) > 0 {
										t := field.Names[0]
										i.accessor = t.Name
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return i
}

var packageTemplate = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
// This file was generated by evoli at
// {{ .Timestamp }}
package main

import (
	"embed"
	{{- range .Imports}}
	"{{ . }}"
	{{- end }}
)

//go:embed resources public configs
var content embed.FS

func (app *{{.App}}) Start() {
    app.{{.Accessor}} = evoli.NewApplication()
	app.{{.Accessor}}.AddEmbedFS(content)

	{{- range $func, $route := .Routes}}
	app.{{$.Accessor}}.AddRoutes("{{$route}}", {{$func}})
	{{- end }}

	{{- range .Components}}
	app.{{$.Accessor}}.RegisterComponent({{ . }})
	{{- end }}

	app.{{.Accessor}}.Start()
}

`))

func generate() {
	command := "generate"
	definition := "generate"
	def := console.Parse(definition, command)

	cmd := Generate()
	cmd.Execution(def)
}
