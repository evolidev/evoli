package view

type Engine struct {
	RenderData map[string]any
}

func NewEngine() *Engine {
	return &Engine{
		RenderData: make(map[string]any),
	}
}

func (e *Engine) AddRenderData(key string, data any) {
	e.RenderData[key] = data
}
