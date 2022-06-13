package main

import (
	"github.com/evolidev/evoli/framework/console/reload"
	"github.com/evolidev/evoli/framework/logging"
	"log"
)

var (
	env  *string
	port *int
)

var logger = logging.NewLogger(&logging.Config{Name: "cli"})

func init() {

}

func main() {
	// Create new parser object
	//parser := argparse.NewParser("print", "Prints provided string to stdout")
	//// Create string flag
	////s := parser.String("s", "string", &argparse.Options{Required: true, Help: "String to print"})
	//// Parse input
	//err := parser.Parse(os.Args)
	//if err != nil {
	//	// In case of error print error and print usage
	//	// This can also be done by passing -h or --help flags
	//	fmt.Print(parser.Usage(err))
	//}
	// Finally print the collected string
	//fmt.Println(*s)

	//console.Commands()
	Watch()
}

func Watch() {
	config := &reload.Configuration{
		AppRoot:            "/Users/omohamed/Code/demo2/evoli/examples/simple",
		IncludedExtensions: []string{".go"},
		BuildPath:          "",
		BinaryName:         "main.go",
		Command:            "go run main.go",
		Debug:              false,
		ForcePolling:       false,
	}
	r := reload.RunBackground(config)
	log.Println(r)
}
