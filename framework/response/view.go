package response

import (
	"bytes"
	"github.com/evolidev/evoli/framework/config"
	"github.com/evolidev/evoli/framework/use"
	"github.com/evolidev/evoli/framework/view"
	"net/http"
	"path"
	"strings"
	"text/template"
)

type ViewResponse struct {
	baseResponse
	template string
	data     map[string]interface{}
	layout   string
	config   *config.Config
}

func View(template string) *ViewResponse {
	response := &ViewResponse{
		config:       use.Config("view"),
		template:     template,
		baseResponse: baseResponse{myHeaders: use.NewCollection[string, string]()},
	}

	response.ensureConfigSettings()

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
	tmpl := r.parseTemplate()

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
	r.layout = layout

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

func (r *ViewResponse) ensureConfigSettings() {
	r.config.SetDefault("path", "resources/views")
	r.config.SetDefault("delimiters.left", "${")
	r.config.SetDefault("delimiters.right", "}")
	r.config.SetDefault("extension", "html")
}

func (r *ViewResponse) parseTemplate() *template.Template {
	files := r.getFilesToParse()

	fs := use.Storage().FS()

	tmpl := template.New(path.Base(files[0]))
	tmpl.Delims(r.config.Get("delimiters.left").Value().(string), r.config.Get("delimiters.right").Value().(string))
	tmpl = template.Must(tmpl.ParseFS(fs, files...))

	return tmpl
}

func (r *ViewResponse) getFilesToParse() []string {
	tmpDir := r.config.Get("path").Value().(string)
	extension := r.config.Get("extension").Value().(string)

	files := make([]string, 0)
	file := strings.Replace(r.template, ".", "/", -1) + "." + extension
	files = append(files, tmpDir+"/"+file)

	if r.layout != "" {
		layout := strings.Replace(r.layout, ".", "/", -1) + ".html"
		files = append(files, tmpDir+"/"+layout)
	}

	return files
}
