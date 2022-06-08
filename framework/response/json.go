package response

import (
	"encoding/json"
	"github.com/evolidev/evoli/framework/use"
)

type JsonResponse struct {
	obj       interface{}
	myHeaders *use.Collection[string, string]
}

func Json(obj interface{}) *JsonResponse {
	return &JsonResponse{obj: obj, myHeaders: use.NewCollection[string, string]()}
}

func (r *JsonResponse) AsBytes() []byte {
	result, err := json.Marshal(r.obj)
	if err != nil {
		panic(err)
	}

	return result
}

func (r *JsonResponse) Headers() *use.Collection[string, string] {
	r.myHeaders.Add("Content-Type", "application/json")

	return r.myHeaders
}
