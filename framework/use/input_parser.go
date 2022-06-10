package use

import (
	"reflect"
	"strconv"
)

type inputParser struct {
	inputType reflect.Type
	value     interface{}
}

func newInputParser(inputType reflect.Type, value interface{}) *inputParser {
	return &inputParser{inputType: inputType, value: value}
}

func (i *inputParser) parse() reflect.Value {
	kind := i.inputType.Kind()

	parsedParam := reflect.ValueOf(i.value)

	if kind == reflect.String {
		parsedParam = i.parseStringParam(i.value)
	} else if kind == reflect.Int {
		parsedParam = i.parseIntParam(i.value)
	} else if kind == reflect.Bool {
		parsedParam = i.parseBoolParam(i.value)
	}

	return parsedParam
}

func (i *inputParser) parseStringParam(param interface{}) reflect.Value {
	parsedParam := reflect.ValueOf(param)
	if reflect.ValueOf(param).Kind() == reflect.Int {
		parsedParam = reflect.ValueOf(strconv.Itoa(param.(int)))
	}

	return parsedParam
}

func (i *inputParser) parseIntParam(param interface{}) reflect.Value {
	parsedParam := reflect.ValueOf(param)
	if reflect.ValueOf(param).Kind() == reflect.String {
		converted, _ := strconv.Atoi(param.(string))
		parsedParam = reflect.ValueOf(converted)
	}

	return parsedParam
}

func (i *inputParser) parseBoolParam(param interface{}) reflect.Value {
	parsedParam := reflect.ValueOf(param)

	if reflect.ValueOf(param).Kind() == reflect.String {
		converted, _ := strconv.ParseBool(param.(string))
		parsedParam = reflect.ValueOf(converted)
	}

	return parsedParam
}
