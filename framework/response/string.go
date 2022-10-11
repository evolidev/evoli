package response

import (
	"evoli.dev/framework/use"
)

type StringResponse struct {
	baseResponse
	body      string
	myHeaders *use.Collection[string, string]
}

func String(body string) *StringResponse {
	return &StringResponse{body: body, myHeaders: use.NewCollection[string, string]()}
}

func (r *StringResponse) AsBytes() []byte {
	return []byte(r.body)
}

func (r *StringResponse) WithCode(code int) *StringResponse {
	r.code = code

	return r
}

func (r *StringResponse) Headers() *use.Collection[string, string] {
	r.myHeaders.Add("Content-Type", "text/plain")

	return r.myHeaders
}
