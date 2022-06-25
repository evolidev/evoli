package component

import (
	"github.com/evolidev/evoli/framework/use"
	"github.com/evolidev/evoli/framework/view"
)

var components = make(map[string]Component)

type Component interface {
}

type Data struct {
	Name string `json:"name"`
}

type Request struct {
	Component  string                 `json:"component"`
	Method     string                 `json:"method"`
	State      map[string]interface{} `json:"state"`
	Action     string                 `json:"action"`
	Parameters []interface{}          `json:"parameters"`
}

type Response struct {
	Component string                 `json:"component"`
	State     map[string]interface{} `json:"state"`
	Type      string                 `json:"type"`
	Content   string                 `json:"content"`
	Response  interface{}            `json:"response"`
}

func New(componentStruct interface{}, data map[string]interface{}) *Base {
	collection := use.NewCollection[string, interface{}]()
	collection.Set(data)

	component := use.Magic(componentStruct).ToPointer()

	return &Base{
		Component: component.WithParams(data).Fill(),
		Data:      collection,
	}
}

func Register(component Component) {
	name := use.GetInterfacedStructName(component)
	components[name] = component
}

func GetRegisteredComponents() *map[string]Component {
	return &components
}

func GetRegisterComponentsCount() int {
	return len(components)
}

func NewByNameWithData(name string, data string) *Base {
	componentObject, ok := components[name]

	if !ok {
		return nil
	}

	mappedData := use.JsonDecodeObject(data)
	component := New(componentObject, mappedData)

	return component
}

func NewByName(name string, data map[string]interface{}) *Base {
	componentObject, ok := components[name]

	if !ok {
		return nil
	}

	component := New(componentObject, data)

	return component
}

func Handle(request *Request) *Response {
	component := NewByName(request.Component, request.State)

	var response interface{}

	if request.Action == "click" {
		response = component.Call(request.Method, request.Parameters)
	}

	return &Response{
		Component: request.Component,
		State:     component.GetState(),
		Response:  response,
	}
}

type Methods struct{}

func (c *Methods) Include(name string) string {
	use.D("include component file: " + name)
	return NewByName(name, nil).Render()
}

func SetupViewEngine(engine *view.Engine) {
	engine.AddRenderData("Component", &Methods{})
	engine.AddPlaceholder("@componentHeader", `
	<script src="https://cdn.tailwindcss.com"></script>
`)
	engine.AddPlaceholder("@componentFooter", `<!-- @componentFooter -->`)
}
