package web

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/markbates/refresh/refresh"
)

var lpath = refresh.ErrorLogPath()
var tmpl *template.Template

func init() {
	tmpl, _ = template.New("template").Parse(html)
}

func ErrorChecker(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ee, err := ioutil.ReadFile(lpath)
		if err != nil || ee == nil {
			h.ServeHTTP(res, req)
			return
		}
		res.WriteHeader(500)
		tmplErr := tmpl.Execute(res, string(ee))
		if tmplErr != nil {
			// todo log to our logger
			fmt.Println(tmplErr)
		}
	})
}

var html = `
<html>
<head>
	<title>Refresh Build Error!</title>
	<style>
		body {
			margin-top: 20px;
			font-family: Helvetica;
		}
		h1 {
			text-align: center;
		}
		pre {
			border: 1px #B22222 solid;
			background-color: #FFB6C1;
			padding: 5px;
			font-size: 32px;
		}
	</style>
</head>

<h1>Oops!! There was a build error!</h1>

<pre><code>{{.}}</code></pre>

</html>
`
