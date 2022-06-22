package middleware

import (
	"encoding/json"
	"github.com/evolidev/evoli/framework/logging"
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
		logger: logging.NewLogger(&logging.Config{Name: "simple application", PrefixColor: 73}),
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
		jsonResponse, _ := json.Marshal(info)
		lm.logger.Success(jsonResponse)

		next.ServeHTTP(writer, request)
	})
}
