package component

import (
	"github.com/evolidev/evoli/framework/use"
)

var components = make(map[string]Component)

type Component interface {
}

func New(component Component) *Base {
	return &Base{
		Component: component,
		Data:      use.NewCollection[string, interface{}](),
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

	component := New(componentObject)

	component.Set(mappedData)

	return component
}

//func (arcane *Arcane) RegisterComponents(components ...Component) {
//	myComponents = map[string]Component{}
//	arcane.routes = make(map[string]ComponentHandler)
//	for _, component := range components {
//		name := reflect.ValueOf(component).Type().Name()
//
//		myComponents[ToSnakeCase(name)] = component
//
//		if page, ok := component.(PageComponent); ok {
//			routeName := page.Route()
//
//			componentHandler := ComponentHandler{Component: component}
//			fmt.Println("Registering component", name)
//
//			arcane.routes[routeName] = componentHandler
//		}
//	}
//}
