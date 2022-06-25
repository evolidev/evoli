package view

type Engine struct {
	RenderData   map[string]any
	Placeholders map[string]string
}

func NewEngine() *Engine {
	return &Engine{
		RenderData:   make(map[string]any),
		Placeholders: make(map[string]string),
	}
}

func (e *Engine) AddRenderData(key string, data any) {
	e.RenderData[key] = data
}

func (e *Engine) AddPlaceholder(key string, data string) {
	e.Placeholders[key] = data
}
