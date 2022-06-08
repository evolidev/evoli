package response

import (
	"github.com/evolidev/evoli/framework/use"
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
	default:
		return Json(arg)
	}
}

type Response interface {
	AsBytes() []byte
	Headers() *use.Collection[string, string]
}
