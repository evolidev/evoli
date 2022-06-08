package component

import (
	"github.com/evolidev/evoli/framework/use"
)

var components = make(map[string]Component)

type Component interface {
}

func New(componentStruct interface{}, data map[string]interface{}) *Base {

	//log.Println("INIT NEW")
	//use.HasMethod(componentStruct, "Test")
	//use.HasMethod(&componentStruct, "Test")
	//use.HasMethod(component, "Test")
	//use.HasMethod(component, "Test")
	//use.HasMethod(&component, "Test")
	//
	//if data != nil {
	//	dataJson := use.JsonEncode(data)
	//	if err := json.Unmarshal([]byte(dataJson), component); err != nil {
	//		// TODO log error
	//		log.Println("Unable to parse json")
	//	}
	//}

	return &Base{
		Component: use.Magic(componentStruct),
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

	component := New(componentObject, nil)

	mappedData := use.JsonDecodeObject(data)
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
