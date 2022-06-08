package use

import (
	"reflect"
	"strconv"
)

func Magic(i interface{}) *Reflection {
	return &Reflection{
		t: reflect.TypeOf(i),
		v: reflect.ValueOf(i),
		p: NewCollection[string, reflect.Value](),
		m: "",
	}
}

type Reflection struct {
	t reflect.Type
	v reflect.Value
	p *Collection[string, reflect.Value]
	m string
}

func (r *Reflection) Call() reflect.Value {
	var result []reflect.Value

	if r.m != "" {
		result = r.v.MethodByName(r.m).Call(r.p.Slice())
	} else {
		result = r.v.Call(r.p.Slice())
	}

	return result[0]
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
	case map[string]string:
		for key, value := range params.(map[string]string) {
			r.p.Add(key, reflect.ValueOf(value))
		}
	}

	return r
}

func (r *Reflection) Method(method string) *Reflection {
	r.m = method

	return r
}
