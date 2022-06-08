package use

import (
	"log"
	"reflect"
)

func GetInterfacedStructName(element interface{}) string {
	return reflect.ValueOf(element).Type().Name()
}

func HasMethod(element interface{}, methodName string) (bool, reflect.Value) {
	//data := reflect.ValueOf(element)
	//P(reflect.New(reflectElem(element)).Elem())

	data := reflect.New(reflectElem(element)).Elem()
	method := data.MethodByName(methodName)

	log.Println("HasMethod", methodName, method.IsValid(), method)

	return method.IsValid(), method
}

func reflectElem(element interface{}) reflect.Type {
	reflectType := reflect.TypeOf(element)
	if reflectType.Kind() == reflect.Ptr {
		return reflectType.Elem()
	}

	return reflectType
}
