package middleware

import (
	"github.com/evolidev/evoli/framework/console/color"
	"github.com/evolidev/evoli/framework/logging"
	"github.com/evolidev/evoli/framework/use"
	"io/ioutil"
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
		timer := use.TimeRecord()

		err := request.ParseForm()
		if err != nil {
			lm.logger.Error(err)
		}
		body := make([]byte, 0)
		if request.Body != nil {
			var err error
			body, err = ioutil.ReadAll(request.Body)
			if err != nil {
				lm.logger.Error(err)
			}
		}

		info := requestInfo{
			Uri:    request.RequestURI,
			Body:   string(body),
			Form:   request.Form,
			Method: request.Method,
		}

		next.ServeHTTP(writer, request)

		jsonResponse := use.JsonEncode(info)

		lm.logger.Log("%s %s", jsonResponse, color.Text(150, timer.ElapsedColored()))
	})
}
