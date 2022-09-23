package use

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mitchellh/mapstructure"
	"net/url"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

func Magic(i interface{}) *Reflection {
	return &Reflection{
		t:       reflect.TypeOf(i),
		v:       reflect.ValueOf(i),
		p:       NewCollection[string, interface{}](),
		injects: NewCollection[string, interface{}](),
		ptr:     reflect.New(reflect.TypeOf(i)),
	}
}

type Reflection struct {
	t       reflect.Type
	v       reflect.Value
	p       *Collection[string, interface{}]
	ptr     reflect.Value
	injects *Collection[string, interface{}]
}

func (r *Reflection) Call() reflect.Value {
	result := r.v.Call(r.parseParams())

	if len(result) > 0 {
		return result[0]
	}

	return reflect.ValueOf(nil)
}

func (r *Reflection) Fill() *Reflection {
	if r.t.Kind() == reflect.Func {
		attributes := make([]reflect.Value, 0)
		arguments := r.appendReceiver(attributes)
		if len(arguments) > 0 {
			receiver := arguments[0]

			return getDecodedDestination(receiver, r.p).Method(r.Name()).WithInjectable(r.p.Slice())
		}

		r.WithInjectable(r.p.Slice())
		r.p = NewCollection[string, interface{}]()

		return r
	}

	return getDecodedDestination(reflect.New(r.reflectElem()), r.p)
}

func getDecodedDestination(input reflect.Value, params *Collection[string, interface{}]) *Reflection {
	destination := getDestination(input).Interface()
	decode(params, destination)

	return Magic(destination)
}

func getDestination(value reflect.Value) reflect.Value {
	destination := value.Interface()
	reflectValue := reflect.ValueOf(destination)

	t := reflectValue.Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return reflect.New(t)
}

func decode(input *Collection[string, interface{}], output interface{}) {

	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   output,
		Squash:   true,
	}

	decoder, _ := mapstructure.NewDecoder(config)

	err := decoder.Decode(input.Map())

	if err != nil {
		panic(err)
	}
}

func (r *Reflection) ToPointer() *Reflection {
	if r.v.Kind() == reflect.Ptr {
		return r
	}

	p := reflect.New(r.t)
	p.Elem().Set(r.v)

	r.v = p

	return r
}

func (r *Reflection) parseParams() []reflect.Value {
	parser := newParamParser(r)

	return parser.parse()
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

func (r *Reflection) Value() interface{} {
	return r.v.Interface()
}

func (r *Reflection) WithParams(params interface{}) *Reflection {
	switch params.(type) {
	case []string:
		for key, value := range params.([]string) {
			r.p.Add(strconv.Itoa(key), value)
		}
	case []interface{}:
		for key, value := range params.([]interface{}) {
			r.p.Add(strconv.Itoa(key), value)
		}
	case map[string]interface{}:
		for key, value := range params.(map[string]interface{}) {
			r.p.Add(key, value)
		}
	case url.Values:
		data := make(map[string]interface{})
		myParams := params.(url.Values)
		for key, _ := range myParams {
			data[key] = myParams.Get(key)
		}
		r.p.Add("form", data)
	case httprouter.Params:
		for _, param := range params.(httprouter.Params) {
			r.p.Add(param.Key, param.Value)
		}
	case *Collection[string, interface{}]:
		r.p.Merge(params.(*Collection[string, interface{}]))
	}

	return r
}

func (r *Reflection) NewPointer() reflect.Value {
	return reflect.New(r.reflectElem())
}

func (r *Reflection) Method(method string) *Reflection {
	var caller reflect.Value

	if r.v.MethodByName(method).IsValid() {
		caller = r.v.MethodByName(method)
	} else if r.ptr.MethodByName(method).IsValid() {
		caller = r.ptr.MethodByName(method)
	}

	return &Reflection{
		v:       caller,
		t:       caller.Type(),
		p:       NewCollection[string, interface{}](),
		ptr:     reflect.New(caller.Type()),
		injects: NewCollection[string, interface{}](),
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

func (r *Reflection) GetField(field string) interface{} {
	return reflect.Indirect(r.v).FieldByName(field).Interface()
}

func (r *Reflection) GetFields() map[string]interface{} {
	fields := make(map[string]interface{})
	indirect := reflect.Indirect(r.v)
	count := indirect.NumField()

	for i := 0; i < count; i++ {
		field := indirect.Type().Field(i)
		fields[field.Name] = indirect.Field(i).Interface()
	}

	return fields
}

func (r *Reflection) WithInjectable(params []interface{}) *Reflection {
	for _, value := range params {
		r.injects.Add(reflect.TypeOf(value).String(), value)
	}

	return r
}
