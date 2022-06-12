package main

import (
	"flag"
	"fmt"
	"github.com/evolidev/evoli/framework/console"
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

	console.Commands()
}

func Watch() {
	config := &reload.Configuration{
		AppRoot:            "/Users/omer/Code/evoli/examples/simple",
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
