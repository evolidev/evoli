package use

import (
	"reflect"
	"strconv"
)

type inputParser struct {
	method *Reflection
}

func newParamParser(method *Reflection) *inputParser {
	return &inputParser{method: method}
}

func (i *inputParser) parse() []reflect.Value {
	parsedParams := make([]reflect.Value, 0)

	cnt := 0
	for cnt < i.method.t.NumIn() {
		currentInputParam := i.method.t.In(cnt)
		cnt++

		inputKind := currentInputParam.String()
		if i.method.injects.Has(inputKind) {
			parsedParams = append(parsedParams, reflect.ValueOf(i.method.injects.Get(inputKind)))

			continue
		}

		param := i.method.p.Next()

		parsedParams = append(parsedParams, i.parseSingle(currentInputParam, param))
	}

	return parsedParams
}

func (i *inputParser) parseSingle(input reflect.Type, value interface{}) reflect.Value {
	kind := input.Kind()

	parsedParam := reflect.ValueOf(value)

	if kind == reflect.String {
		parsedParam = i.parseStringParam(value)
	} else if kind == reflect.Int {
		parsedParam = i.parseIntParam(value)
	} else if kind == reflect.Bool {
		parsedParam = i.parseBoolParam(value)
	}

	return parsedParam
}

func (i *inputParser) appendReceiver(arguments []reflect.Value) []reflect.Value {
	amount := i.method.t.NumIn()

	if amount > 0 {
		first := i.method.t.In(0)

		isPointer := first.Kind() == reflect.Ptr
		if first.Kind() == reflect.Struct || isPointer {
			if _, ok := first.MethodByName(i.method.Name()); ok {
				newReceiver := reflect.New(first)

				if !isPointer {
					newReceiver = newReceiver.Elem()
				}

				arguments = append(arguments, newReceiver)
			}
		}
	}

	return arguments
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
