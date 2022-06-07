package component

import (
	"fmt"
	"github.com/evolidev/evoli/framework/filesystem"
	"github.com/evolidev/evoli/framework/use"
	"reflect"
)

var components map[string]Component

type Component interface {
}

type Base struct {
	File               string
	Component          *Component
	ComponentInterface interface{}
	Data               *use.Collection[string, interface{}]
}

func New(component Component) *Base {
	return &Base{
		Component: &component,
		Data:      use.NewCollection[string, interface{}](),
	}
}

func Register(component Component) {
	name := use.GetInterfacedStructName(component)

	components[name] = component
}

func (b *Base) GetComponentInterface() interface{} {
	if b.ComponentInterface != nil {
		return b.ComponentInterface
	}

	componentInterface := reflect.New(reflect.TypeOf(*b.Component)).Interface()

	b.ComponentInterface = componentInterface

	return b.ComponentInterface
}

func (b *Base) GetFilePath() string {
	if ok, method := use.HasMethod(b.Component, "GetFilePath"); ok {
		output := method.Call([]reflect.Value{})
		return output[0].String()
	}

	return fmt.Sprintf("templates/%s.html", use.String(b.GetComponentName()).Kebab().Get())
}

func (b *Base) GetComponentName() string {
	return use.GetInterfacedStructName(*b.Component)
}

func (b *Base) GetRawContent() string {
	path := b.GetFilePath()
	return filesystem.Read(path)
}

func (b *Base) Render() string {
	content := b.GetRawContent()
	return content
}

func (b *Base) SetData(data map[string]interface{}) {
	b.Data.Set(data)
}

func (b *Base) Get(key string) interface{} {
	return b.Data.Get(key)
}

//func NewByName(name string, data string) *Component {
//	// check if component is registered
//	component, ok := components[name]
//
//	if !ok {
//		return nil
//	}
//
//	componentInterface := reflect.New(reflect.TypeOf(component)).Interface()
//	componentInterface.(Component).SetData(data)
//
//	return componentInterface.(Component)
//}

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
