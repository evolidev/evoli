package response

type StringResponse struct {
	body string
}

func String(body string) *StringResponse {
	return &StringResponse{body: body}
}

func (r *StringResponse) AsBytes() []byte {
	return []byte(r.body)
}

func (r *StringResponse) Headers() map[string]string {
	headers := make(map[string]string)
	headers["Content-Type"] = "text/plain"

	return headers
}
