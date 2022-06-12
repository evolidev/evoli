package main

import (
	"github.com/evolidev/evoli/framework/component"
	"github.com/evolidev/evoli/framework/use"
	"log"
)

func main() {
	log.Println("Starting.")

	consoleTest()
	consoleTest()
	consoleTest()
	consoleTest()

	helloWorldComponentTest()
}

func consoleTest() {
	log.Print("Console test..D")
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
	use.D("It is working!")
	return "super."
}

func helloWorldComponentTest() {
	component.Register(HelloWorldWithPath{})

	json := `{"Name":"Foo"}`

	hello := component.NewByNameWithData("HelloWorldWithPath", json)

	response := hello.Call("Test", nil)

	use.D(response)
}
