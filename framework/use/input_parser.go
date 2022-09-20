package use

import (
	"github.com/mitchellh/mapstructure"
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
	parsedParams = i.appendReceiver(parsedParams)

	cnt := 0
	for cnt < i.method.t.NumIn() {
		currentInputParam := i.method.t.In(cnt)
		if cnt == 0 && len(parsedParams) > 0 {
			cnt++
			continue
		}
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
	} else if (kind == reflect.Struct || kind == reflect.Ptr) && (parsedParam.Kind() == reflect.String || parsedParam.Kind() == reflect.Int) {
		id := i.parseIntParam(value)

		parsedParam = getRecordFromTypeAndId(input, int(id.Int()))
	} else if kind == reflect.Struct {
		form := i.parseFormValues(value)

		destination := reflect.New(input).Elem().Interface()
		reflectValue := reflect.ValueOf(destination)
		destination = reflect.New(reflectValue.Type()).Interface()

		mapstructure.Decode(form, destination)

		parsedParam = reflect.ValueOf(destination).Elem()
	}

	return parsedParam
}

func getRecordFromTypeAndId(paramType reflect.Type, id int) reflect.Value {

	destination := reflect.New(paramType).Elem().Interface()
	reflectValue := reflect.ValueOf(destination)
	destination = reflect.New(reflectValue.Type().Elem()).Interface()

	DB().First(destination, id)

	return reflect.ValueOf(destination)
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

func (i *inputParser) parseFormValues(param interface{}) interface{} {
	return param
}
