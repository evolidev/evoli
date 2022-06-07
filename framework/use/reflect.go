package use

import "reflect"

func GetInterfacedStructName(element interface{}) string {
	return reflect.ValueOf(element).Type().Name()
}

func HasMethod(element interface{}, methodName string) (bool, reflect.Value) {
	data := reflect.ValueOf(element)
	method := data.MethodByName(methodName)

	return method.IsValid(), method
}
