package component

import (
	"github.com/evolidev/evoli/framework/use"
)

var components = make(map[string]Component)

type Component interface {
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
