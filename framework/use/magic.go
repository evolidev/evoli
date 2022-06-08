package use

import (
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

func Magic(i interface{}) *Reflection {
	return &Reflection{
		t:   reflect.TypeOf(i),
		v:   reflect.ValueOf(i),
		p:   NewCollection[string, reflect.Value](),
		ptr: reflect.New(reflect.TypeOf(i)),
	}
}

type Reflection struct {
	t   reflect.Type
	v   reflect.Value
	p   *Collection[string, reflect.Value]
	ptr reflect.Value
}

func (r *Reflection) Call() reflect.Value {
	result := r.v.Call(r.parseParams())

	if len(result) > 0 {
		return result[0]
	}

	return reflect.Value{}
}

func (r *Reflection) parseParams() []reflect.Value {
	var parsedArguments = make([]reflect.Value, 0)

	parsedArguments = r.appendReceiver(parsedArguments)

	return parsedArguments
}

func (r *Reflection) appendReceiver(arguments []reflect.Value) []reflect.Value {
	amount := r.t.NumIn()

	if amount > 0 {
		first := r.t.In(0)

		isPointer := first.Kind() == reflect.Ptr
		if first.Kind() == reflect.Struct || isPointer {
			if _, ok := first.MethodByName(r.Name()); ok {
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

func (r *Reflection) WithParams(params interface{}) *Reflection {
	switch params.(type) {
	case []string:
		for key, value := range params.([]string) {
			r.p.Add(strconv.Itoa(key), reflect.ValueOf(value))
		}
	case []int:
		for key, value := range params.([]int) {
			r.p.Add(strconv.Itoa(key), reflect.ValueOf(value))
		}
	case []interface{}:
		for key, value := range params.([]interface{}) {
			r.p.Add(strconv.Itoa(key), reflect.ValueOf(value))
		}
	case map[string]string:
		for key, value := range params.(map[string]string) {
			r.p.Add(key, reflect.ValueOf(value))
		}
	case map[string]interface{}:
		for key, value := range params.(map[string]interface{}) {
			r.p.Add(key, reflect.ValueOf(value))
		}
	}

	return r
}

func (r *Reflection) NewPointer() reflect.Value {
	return reflect.New(r.t)
}

func (r *Reflection) Method(method string) *Reflection {
	var caller reflect.Value

	if r.v.MethodByName(method).IsValid() {
		caller = r.v.MethodByName(method)
	} else if r.ptr.MethodByName(method).IsValid() {
		caller = r.ptr.MethodByName(method)
	}

	return &Reflection{
		v:   caller,
		t:   caller.Type(),
		p:   NewCollection[string, reflect.Value](),
		ptr: reflect.New(caller.Type()),
	}
}

func (r *Reflection) HasMethod(method string) bool {
	isValid := r.v.MethodByName(method).IsValid()

	if isValid == false {
		isValid = reflect.New(r.t).MethodByName(method).IsValid()
	}

	return isValid
}

func (r *Reflection) Name() string {
	if r.t.Kind() == reflect.Func {
		funcName := r.functionName()
		paths := strings.Split(funcName, ".")

		return paths[len(paths)-1]
	}

	return r.reflectElem().Name()
}

func (r *Reflection) functionName() string {
	return runtime.FuncForPC(r.v.Pointer()).Name()
}

func (r *Reflection) reflectElem() reflect.Type {
	if r.t.Kind() == reflect.Ptr {
		return r.t.Elem()
	}

	return r.t
}