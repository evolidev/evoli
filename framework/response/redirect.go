package response

import (
	"evoli.dev/framework/use"
	"net/http"
)

type RedirectResponse struct {
	baseResponse
	To string
}

func Redirect(to string) *RedirectResponse {
	response := &RedirectResponse{To: to, baseResponse: baseResponse{myHeaders: use.NewCollection[string, string]()}}

	return response.WithCode(http.StatusSeeOther)
}

func (r *RedirectResponse) AsBytes() []byte {
	return []byte{}
}

func (r *RedirectResponse) WithCode(code int) *RedirectResponse {
	r.code = code

	return r
}
