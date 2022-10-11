package response

import (
	"evoli.dev/framework/use"
	"net/http"
)

type EmptyResponse struct {
	baseResponse
}

func Empty() *EmptyResponse {
	response := &EmptyResponse{baseResponse: baseResponse{myHeaders: use.NewCollection[string, string]()}}

	return response.WithCode(http.StatusNoContent)
}

func (r *EmptyResponse) AsBytes() []byte {
	return []byte{}
}

func (r *EmptyResponse) WithHeader(key string, value string) *EmptyResponse {
	r.myHeaders.Add(key, value)

	return r
}

func (r *EmptyResponse) WithCode(code int) *EmptyResponse {
	r.code = code

	return r
}
