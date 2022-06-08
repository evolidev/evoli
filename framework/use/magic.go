package use

import (
	"reflect"
	"strconv"
)

func Magic(i interface{}) *Reflection {
	return &Reflection{
		t:   reflect.TypeOf(i),
		v:   reflect.ValueOf(i),
		p:   NewCollection[string, reflect.Value](),
		m:   "",
		ptr: reflect.New(reflect.TypeOf(i)),
	}
}

type Reflection struct {
	t   reflect.Type
	v   reflect.Value
	p   *Collection[string, reflect.Value]
	m   string
	ptr reflect.Value
}

func (r *Reflection) Call() reflect.Value {
	var result []reflect.Value

	var caller reflect.Value

	if r.m != "" {
		if r.v.MethodByName(r.m).IsValid() {
			caller = r.v.MethodByName(r.m)
		} else if r.ptr.MethodByName(r.m).IsValid() {
			caller = r.ptr.MethodByName(r.m)
		}
	} else {
		caller = r.v
	}

	amount := caller.Type().NumIn()
	arguments := r.p.Slice()

	if amount > 0 {
		first := caller.Type().In(0)

		isPointer := first.Kind() == reflect.Ptr
		if first.Kind() == reflect.Struct || isPointer {
			if _, ok := first.MethodByName(r.Name()); ok {
				newReceiver := reflect.New(first)

				if !isPointer {
					newReceiver = newReceiver.Elem()
				}

				arguments = append([]reflect.Value{newReceiver}, arguments...)
			}
		}
	}

	cnt := 0

	var parsedArguments = make([]reflect.Value, 0)
	for _, value := range arguments {
		kind := caller.Type().In(cnt).Kind()

		if kind == reflect.Int && value.Type().Kind() == reflect.String {
			parsedParam, _ := strconv.Atoi(value.Interface().(string))
			parsedArguments = append(parsedArguments, reflect.ValueOf(parsedParam))
		} else if kind == reflect.Int {
			// TODO
		} else {
			parsedArguments = append(parsedArguments, reflect.ValueOf(value))
		}

		cnt++
	}

	result = caller.Call(parsedArguments)

	if len(result) > 0 {
		return result[0]
	}

	return reflect.Value{}
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

func (r *Reflection) Method(method string) *Reflection {
	r.m = method

	return r
}

func (r *Reflection) HasMethod(method string) bool {
	isValid := r.v.MethodByName(method).IsValid()

	if isValid == false {
		isValid = reflect.New(r.t).MethodByName(method).IsValid()
	}

	return isValid
}

func (r *Reflection) Name() string {
	return r.reflectElem().Name()
}

func (r *Reflection) reflectElem() reflect.Type {
	if r.t.Kind() == reflect.Ptr {
		return r.t.Elem()
	}

	return r.t
}
