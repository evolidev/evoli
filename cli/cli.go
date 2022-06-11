package main

import (
	"github.com/evolidev/evoli/framework/console/reload"
	"log"
)

func main() {
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
