package response

import "encoding/json"

type JsonResponse struct {
	obj interface{}
}

func Json(obj interface{}) *JsonResponse {
	return &JsonResponse{obj: obj}
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
