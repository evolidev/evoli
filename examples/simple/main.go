package main

import (
	"github.com/evolidev/evoli"
	"github.com/evolidev/evoli/framework/use"
)

//go:generate go run main.go generate

type TestApp struct {
	*evoli.Application
	test string
}

func main() {
	app := &TestApp{}

	use.D(use.BasePath("storage"))

	app.Start()
}
