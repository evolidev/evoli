package evoli

import (
	"encoding/json"
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
	Headers() map[string]string
}

type StringResponse struct {
	body string
}

func (r *StringResponse) AsBytes() []byte {
	return []byte(r.body)
}

func (r *StringResponse) Headers() map[string]string {
	headers := make(map[string]string)
	headers["Content-Type"] = "text/plain"

	return headers
}

type JsonResponse struct {
	obj interface{}
}

func (r *JsonResponse) AsBytes() []byte {
	result, err := json.Marshal(r.obj)
	if err != nil {
		panic(err)
	}

	return result
}

func (r *JsonResponse) Headers() map[string]string {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	return headers
}

func String(body string) *StringResponse {
	return &StringResponse{body: body}
}

func Json(obj interface{}) *JsonResponse {
	return &JsonResponse{obj: obj}
}
