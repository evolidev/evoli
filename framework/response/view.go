package response

import (
	"bytes"
	"github.com/evolidev/evoli/framework/use"
	"html/template"
	"path/filepath"
	"strings"
)

type ViewResponse struct {
	template  string
	myHeaders *use.Collection[string, string]
	data      interface{}
	layout    string
}

func View(template string) *ViewResponse {
	return &ViewResponse{template: template, myHeaders: use.NewCollection[string, string]()}
}

func (r *ViewResponse) AsBytes() []byte {
	view, _ := filepath.Abs(strings.Replace(r.template, ".", "/", -1) + ".html")
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

func (r *ViewResponse) Headers() *use.Collection[string, string] {
	r.myHeaders.Add("Content-Type", "text/html")

	return r.myHeaders
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
	r.layout = layout

	return r
}
