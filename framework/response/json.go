package response

import (
	"encoding/json"
	"github.com/evolidev/evoli/framework/use"
)

type JsonResponse struct {
	baseResponse
	obj interface{}
}

func Json(obj interface{}) *JsonResponse {
	response := &JsonResponse{obj: obj, baseResponse: baseResponse{myHeaders: use.NewCollection[string, string]()}}

	return response.WithHeader("Content-Type", "application/json")
}

func (r *JsonResponse) AsBytes() []byte {
	result, err := json.Marshal(r.obj)
	if err != nil {
		panic(err)
	}

	return result
}

func (r *JsonResponse) WithHeader(key string, value string) *JsonResponse {
	r.myHeaders.Add(key, value)

	return r
}

func (r *JsonResponse) WithCode(code int) *JsonResponse {
	r.code = code

	return r
}
