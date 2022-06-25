package component

import (
	"github.com/evolidev/evoli/framework/use"
)

var components = make(map[string]Component)

type Component interface {
}

type Data struct {
	Name string `json:"name"`
}

type Request struct {
	Component  string         `json:"component"`
	Method     string         `json:"method"`
	State      map[string]any `json:"state"`
	Action     string         `json:"action"`
	Parameters []any          `json:"parameters"`
}

type Response struct {
	Component string         `json:"component"`
	State     map[string]any `json:"state"`
	Type      string         `json:"type"`
	Content   string         `json:"content"`
	Response  any            `json:"response"`
}

func New(componentStruct any, data map[string]any) *Base {
	collection := use.NewCollection[string, any]()
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

func NewByName(name string, data map[string]any) *Base {
	componentObject, ok := components[name]

	if !ok {
		return nil
	}

	component := New(componentObject, data)

	return component
}

func Handle(request *Request) *Response {
	component := NewByName(request.Component, request.State)

	var response any

	if request.Action == "click" {
		response = component.Call(request.Method, request.Parameters)
	}

	return &Response{
		Component: request.Component,
		State:     component.GetState(),
		Response:  response,
	}
}
