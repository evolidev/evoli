package response

import (
	"github.com/evolidev/evoli/framework/use"
	"net/http"
)

type StringResponse struct {
	body      string
	myHeaders *use.Collection[string, string]
}

func String(body string) *StringResponse {
	return &StringResponse{body: body, myHeaders: use.NewCollection[string, string]()}
}

func (r *StringResponse) AsBytes() []byte {
	return []byte(r.body)
}

func (r *StringResponse) Headers() *use.Collection[string, string] {
	r.myHeaders.Add("Content-Type", "text/plain")

	return r.myHeaders
}

func (r *StringResponse) Code() int {
	return http.StatusOK
}
