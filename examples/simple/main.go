package main

import (
	"github.com/evolidev/evoli/framework/component"
	"github.com/evolidev/evoli/framework/logging"
)

var logger *logging.Logger

func main() {
	logger = logging.NewLogger(&logging.Config{Name: "simple", PrefixColor: 73})

	logger.Success("Starting..")

	consoleTest()
	consoleTest()
	consoleTest()
	consoleTest()

	helloWorldComponentTest()
}

func consoleTest() {
	logger.Debug("Console test..D")
}

/**
 * HelloWorldWithPath component
 */

type HelloWorldWithPath struct {
}

func (h *HelloWorldWithPath) GetFilePath() string {
	return "hello-worlds"
}

func (h *HelloWorldWithPath) Test() string {
	logger.Debug("It is working!")
	return "super."
}

func helloWorldComponentTest() {
	component.Register(HelloWorldWithPath{})

	json := `{"Name":"Foo"}`

	hello := component.NewByNameWithData("HelloWorldWithPath", json)

	response := hello.Call("Test", nil)

	logger.Success(response)
}
