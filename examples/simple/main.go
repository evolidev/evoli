package main

import (
	"github.com/evolidev/evoli/framework/component"
	"github.com/evolidev/evoli/framework/use"
	"log"
)

func main() {
	log.Println("Starting and update updated :D...")

	consoleTest()

	helloWorldComponentTest()
}

func consoleTest() {
}

/**
 * HelloWorldWithPath component
 */

type HelloWorldWithPath struct {
}

func (h *HelloWorldWithPath) GetFilePath() string {
	return "hello-world"
}

func (h *HelloWorldWithPath) Test() string {
	use.D("yayayaya")

	return "super"
}

func helloWorldComponentTest() {
	component.Register(HelloWorldWithPath{})

	json := `{"Name":"Foo"}`
	hello := component.NewByNameWithData("HelloWorldWithPath", json)

	response := hello.Call("Test", nil)

	use.D(response)
}
