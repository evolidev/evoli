package test

import (
	"github.com/evolidev/evoli/framework/component"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

type helloWorld struct {
	Name string
}

type helloWorldWithPath struct {
}

func (h *helloWorldWithPath) GetFilePath() string {
	return "not-hello-world"
}

func (h *helloWorldWithPath) TestMethod() string {
	return "hello-world-returned"
}

func (h *helloWorldWithPath) TestMethodWithParameters(number int, value string) string {
	return strconv.Itoa(number) + " " + value
}

func TestRenderCorrectComponent(t *testing.T) {

	t.Run("Return the correct given component filesystem", func(t *testing.T) {
		hello := component.New(helloWorld{}, nil)

		assert.Equal(t, "templates/hello-world.html", hello.GetFilePath())
	})

	//t.Run("Return the correct given component filesystem", func(t *testing.T) {
	//	hello := component.New(helloWorld{}, nil)
	//
	//	assert.Equal(t, "templates/hello-world.html", hello.GetFilePath())
	//})

	t.Run("Get the component filesystem content", func(t *testing.T) {
		hello := component.New(helloWorld{}, nil)

		assert.Equal(t, "<div>Hello {{ Name }}</div>", hello.GetRawContent())
	})
	//
	t.Run("Return component with a given path", func(t *testing.T) {
		hello := component.New(helloWorldWithPath{}, nil)

		assert.Equal(t, "not-hello-world", hello.GetFilePath())
	})
	//
	t.Run("Render component with Json data", func(t *testing.T) {
		hello := component.New(helloWorldWithPath{}, nil)
		hello.Set(map[string]interface{}{"Name": "Super"})

		assert.Equal(t, "Super", hello.Get("Name"))
	})
	//
	t.Run("Try to get component that is not registered", func(t *testing.T) {
		json := `{"Name":"Foo"}`
		hello := component.NewByNameWithData("helloWorldWithPath", json)

		assert.Nil(t, hello)
	})
	//
	t.Run("Make sure that components are empty", func(t *testing.T) {
		assert.Equal(t, 0, component.GetRegisterComponentsCount())
	})

	t.Run("Register a component and check if it exists", func(t *testing.T) {
		component.Register(helloWorldWithPath{})

		assert.Equal(t, 1, component.GetRegisterComponentsCount())
	})
	//
	t.Run("Render component with JSON data by name", func(t *testing.T) {
		component.Register(helloWorldWithPath{})

		json := `{"Name":"Foo"}`
		hello := component.NewByNameWithData("helloWorldWithPath", json)
		assert.NotNil(t, hello)

		assert.Equal(t, "Foo", hello.Get("Name"))
	})
	//
	t.Run("Call method of component", func(t *testing.T) {
		component.Register(helloWorldWithPath{})

		hello := component.NewByNameWithData("helloWorldWithPath", `{"Name":"Foo"}`)

		response := hello.Call("TestMethod", nil)

		assert.Equal(t, "hello-world-returned", response.(string))
	})
	//
	t.Run("Call method of component with parameters", func(t *testing.T) {
		hello := component.NewByNameWithData("helloWorldWithPath", `{"Name":"Foo"}`)

		parameters := []interface{}{10, "super"}
		response := hello.Call("TestMethodWithParameters", parameters)

		assert.Equal(t, "10 super", response.(string))
	})

	//t.Run("Call method of component and update property", func(t *testing.T) {
	//	hello := component.NewByNameWithData("helloWorldWithPath", `{"Name":"Foo"}`)
	//
	//	parameters := []interface{}{"FooUpdated"}
	//	hello.Call("UpdateName", parameters)
	//
	//	assert.Equal(t, "FooUpdated", hello.Get("Name"))
	//})
}
