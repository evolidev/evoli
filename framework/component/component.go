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

//type Request struct {
//	Component *Data                  `json:"component"`
//	Method    string                 `json:"method"`
//	State     map[string]interface{} `json:"state"`
//	Action    string                 `json:"action"`
//	Value     string                 `json:"value"`
//}
//
//type Response struct {
//	Component *Data       `json:"component"`
//	State     interface{} `json:"state"`
//	Type      string      `json:"type"`
//	Content   string      `json:"content"`
//}

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
