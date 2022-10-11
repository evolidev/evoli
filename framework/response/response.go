package response

import (
	"evoli.dev/framework/use"
	"reflect"
	"strconv"
)

func NewResponse(arg interface{}) Response {

	switch arg.(type) {
	case Response:
		return arg.(Response)
	case int:
		return String(strconv.Itoa(arg.(int)))
	case string:
		return String(arg.(string))
	case reflect.Value:
		refl := arg.(reflect.Value)
		if !refl.IsValid() || refl.IsZero() {
			return Empty()
		}

		return NewResponse(refl.Interface())
	default:
		return Json(arg)
	}
}

type baseResponse struct {
	myHeaders *use.Collection[string, string]
	code      int
}

func (r *baseResponse) Headers() *use.Collection[string, string] {
	return r.myHeaders
}

func (r *baseResponse) Code() int {
	return r.code
}

type Response interface {
	AsBytes() []byte
	Headers() *use.Collection[string, string]
	Code() int
}
