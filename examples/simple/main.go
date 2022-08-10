package main

import (
	"github.com/evolidev/evoli"
)

//go:generate go run main.go generate

type TestApp struct {
	test string
	*evoli.Application
}

func main() {
	app := &TestApp{}

	app.Start()
}
