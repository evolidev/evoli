package evoli

type StringResponse struct {
	body string
}

type JsonResponse struct {
	obj interface{}
}

func String(body string) *StringResponse {
	return &StringResponse{body: body}
}

func Json(obj interface{}) *JsonResponse {
	return &JsonResponse{obj: obj}
}
