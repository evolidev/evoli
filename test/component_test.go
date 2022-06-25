package test

import (
	"github.com/evolidev/evoli/framework/component"
	"github.com/evolidev/evoli/framework/use"
	"github.com/evolidev/evoli/framework/view"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

type ParentComponent struct {
}

type ChildComponent struct {
}

type helloWorld struct {
	Name string
}

type helloWorldWithPath struct {
	Name string
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

func (h *helloWorldWithPath) UpdateName(value string) {
	h.Name = value
}

func TestRenderCorrectComponent(t *testing.T) {
	t.Parallel()
	t.Run("Return the correct given component filesystem", func(t *testing.T) {
		hello := component.New(helloWorld{}, nil)

		assert.Equal(t, "components/hello-world", hello.GetFilePath())
	})

	t.Run("Return the correct given component filesystem", func(t *testing.T) {
		hello := component.New(helloWorld{}, nil)

		assert.Equal(t, "components/hello-world", hello.GetFilePath())
	})

	t.Run("Get the component filesystem content and make sure the double braces are not parsed", func(t *testing.T) {
		hello := component.New(helloWorld{}, nil)

		assert.Equal(t, "<div>Hello {{ Name }}</div>", hello.GetRawContent())
	})

	t.Run("Return component with a given path", func(t *testing.T) {
		hello := component.New(helloWorldWithPath{}, nil)

		assert.Equal(t, "not-hello-world", hello.GetFilePath())
	})
}

func TestPropertyOfComponents(t *testing.T) {
	t.Run("Do not return a not registered component", func(t *testing.T) {
		hello := component.NewByName("NotExistingComponent", nil)

		assert.Nil(t, hello)
	})

	t.Run("Render component with Json data", func(t *testing.T) {
		hello := component.New(helloWorldWithPath{}, nil)
		hello.Set(map[string]interface{}{"Name": "Super"})

		assert.Equal(t, "Super", hello.Get("Name"))
	})

	t.Run("Try to get component that is not registered", func(t *testing.T) {
		json := `{"Name":"Foo"}`
		hello := component.NewByNameWithData("helloWorldWithPath", json)

		assert.Nil(t, hello)
	})

	t.Run("Make sure that components are empty", func(t *testing.T) {
		assert.Equal(t, 0, component.GetRegisterComponentsCount())
	})

	t.Run("Register a component and check if it exists", func(t *testing.T) {
		component.Register(helloWorldWithPath{})

		assert.Equal(t, 1, component.GetRegisterComponentsCount())
	})

	t.Run("Render component with JSON data by name", func(t *testing.T) {
		component.Register(helloWorldWithPath{})

		json := `{"Name":"Foo"}`
		hello := component.NewByNameWithData("helloWorldWithPath", json)
		assert.NotNil(t, hello)

		assert.Equal(t, "Foo", hello.Get("Name"))
	})

	t.Run("Call method of component", func(t *testing.T) {
		component.Register(helloWorldWithPath{})

		hello := component.NewByNameWithData("helloWorldWithPath", `{"Name":"Foo"}`)

		response := hello.Call("TestMethod", nil)

		assert.Equal(t, "hello-world-returned", response.(string))
	})

	t.Run("Call method of component with parameters", func(t *testing.T) {
		component.Register(helloWorldWithPath{})
		hello := component.NewByNameWithData("helloWorldWithPath", `{"Name":"Foo"}`)

		parameters := []interface{}{10, "super"}
		response := hello.Call("TestMethodWithParameters", parameters)

		assert.Equal(t, "10 super", response.(string))
	})

	t.Run("Call method of component and update property", func(t *testing.T) {
		component.Register(helloWorldWithPath{})
		hello := component.NewByNameWithData("helloWorldWithPath", `{"Name":"Foo"}`)

		assert.Equal(t, "Foo", hello.Get("Name"))

		parameters := []interface{}{"FooUpdated"}
		hello.Call("UpdateName", parameters)

		assert.Equal(t, "FooUpdated", hello.Get("Name"))
	})

}

func TestComponentRequestResponseHandling(t *testing.T) {
	t.Parallel()

	t.Run("Make request to component handler", func(t *testing.T) {
		component.Register(helloWorldWithPath{})

		request := &component.Request{
			Component:  "helloWorldWithPath",
			State:      map[string]interface{}{"Name": "Foo"},
			Action:     "click",
			Method:     "TestMethodWithParameters",
			Parameters: []interface{}{1, "string"},
		}

		response := component.Handle(request)

		assert.Equal(t, response.Response, "1 string")
	})

	t.Run("Make request to and update the name property", func(t *testing.T) {
		component.Register(helloWorldWithPath{})

		request := &component.Request{
			Component:  "helloWorldWithPath",
			State:      map[string]interface{}{"Name": "Foo"},
			Action:     "click",
			Method:     "UpdateName",
			Parameters: []interface{}{"FooUpdated"},
		}

		response := component.Handle(request)

		assert.Equal(t, response.State["Name"], "FooUpdated")
	})
}

func TestComponentRendering(t *testing.T) {
	t.Parallel()

	t.Run("Include component in the page", func(t *testing.T) {
		viewEngine := view.NewEngine()
		component.SetupViewEngine(viewEngine)

		use.AddFacade("viewEngine", viewEngine)

		component.Register(ParentComponent{})
		component.Register(ChildComponent{})

		parent := component.NewByName("ParentComponent", nil)

		content := parent.Render()

		assert.Contains(t, content, "child component")
	})
}
