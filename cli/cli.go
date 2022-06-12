package main

import (
	"flag"
	"fmt"
	"github.com/evolidev/evoli/framework/console/reload"
	"log"
)

var (
	env  *string
	port *int
)

func init() {
	env = flag.String("env", "development", "current environment")
	port = flag.Int("port", 3000, "port number")
}

func main() {
	flag.Parse()

	fmt.Println("env:", *env)
	fmt.Println("port:", *port)

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
