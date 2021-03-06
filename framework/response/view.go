package response

import (
	"bytes"
	"github.com/evolidev/evoli/framework/use"
	"github.com/evolidev/evoli/framework/view"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

type ViewResponse struct {
	baseResponse
	template string
	data     map[string]interface{}
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
	return []byte(r.AsString())
}

func (r *ViewResponse) AsString() string {
	b, err := r.parse()

	if err != nil {
		return ""
	}

	output := b.String()

	return r.replaceTemplate(output)
}

func (r *ViewResponse) parse() (bytes.Buffer, error) {
	rootDir := use.BasePath()
	view, _ := filepath.Abs(rootDir + r.basePath + strings.Replace(r.template, ".", "/", -1) + ".html")

	files := []string{view}

	if r.layout != "" {
		layout, _ := filepath.Abs(strings.Replace(r.layout, ".", "/", -1) + ".html")
		files = append(files, layout)
	}

	tmp := template.New(path.Base(files[0]))
	tmp.Delims("${", "}") // set delimiters (TODO read from config)
	tmpl := template.Must(tmp.ParseFiles(files...))

	var b bytes.Buffer
	err := tmpl.Execute(&b, r.GetAllData())
	return b, err
}

func (r *ViewResponse) WithHeader(key string, value string) *ViewResponse {
	r.myHeaders.Add(key, value)

	return r
}

func (r *ViewResponse) WithData(data map[string]interface{}) *ViewResponse {
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

func (r *ViewResponse) GetAllData() any {
	data := r.data

	if data == nil {
		data = make(map[string]interface{})
	}

	engine := r.GetEngine()

	if engine != nil {

		for key, item := range engine.RenderData {
			data[key] = item
		}
	}

	return data
}

func (r *ViewResponse) GetEngine() *view.Engine {
	engine := use.GetFacade("viewEngine")

	if engine != nil {
		return engine.(*view.Engine)
	}

	return nil
}

func (r *ViewResponse) replaceTemplate(template string) string {
	s := template[:]

	engine := r.GetEngine()

	if engine != nil {
		for key, item := range engine.Placeholders {
			s = strings.ReplaceAll(s, key, item)
		}
	}

	return s
}
