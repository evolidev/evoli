package middleware

import (
	"github.com/evolidev/evoli/framework/logging"
	"github.com/evolidev/evoli/framework/use"
	"net/http"
	"net/url"
)

type requestInfo struct {
	Uri    string
	Body   string
	Form   url.Values
	Method string
}

type LoggingMiddleware struct {
	logger *logging.Logger
}

func NewLoggingMiddleware() LoggingMiddleware {
	return LoggingMiddleware{
		logger: logging.NewLogger(&logging.Config{Name: "router", PrefixColor: 144}),
	}
}

func (lm LoggingMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		body := make([]byte, 0)
		request.Body.Read(body)
		request.ParseForm()
		info := requestInfo{
			Uri:    request.RequestURI,
			Body:   string(body),
			Form:   request.Form,
			Method: request.Method,
		}
		jsonResponse := use.JsonEncode(info)
		lm.logger.Log(jsonResponse)

		next.ServeHTTP(writer, request)
	})
}
