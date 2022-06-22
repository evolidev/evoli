package response

import (
	"bytes"
	"github.com/evolidev/evoli/framework/use"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

type ViewResponse struct {
	baseResponse
	template string
	data     interface{}
	layout   string
	basePath string
}

func View(template string) *ViewResponse {
	response := &ViewResponse{
		template:     template,
		basePath:     "resources/views/", //todo get from config
		baseResponse: baseResponse{myHeaders: use.NewCollection[string, string]()},
	}

	return response.WithCode(http.StatusOK).WithHeader("Content-Type", "text/html")
}

func (r *ViewResponse) AsBytes() []byte {
	view, _ := filepath.Abs(r.basePath + strings.Replace(r.template, ".", "/", -1) + ".html")

	files := []string{view}

	if r.layout != "" {
		layout, _ := filepath.Abs(strings.Replace(r.layout, ".", "/", -1) + ".html")
		files = append(files, layout)
	}

	tmpl := template.Must(template.ParseFiles(files...))

	var b bytes.Buffer
	err := tmpl.Execute(&b, r.data)

	if err != nil {
		return []byte{}
	}

	return b.Bytes()
}

func (r *ViewResponse) WithHeader(key string, value string) *ViewResponse {
	r.myHeaders.Add(key, value)

	return r
}

func (r *ViewResponse) WithData(data interface{}) *ViewResponse {
	r.data = data

	return r
}

func (r *ViewResponse) WithLayout(layout string) *ViewResponse {
	r.layout = r.basePath + "/" + layout

	return r
}

func (r *ViewResponse) SetBasePath(path string) *ViewResponse {
	r.basePath = path

	return r
}

func (r *ViewResponse) WithCode(code int) *ViewResponse {
	r.code = code

	return r
}
