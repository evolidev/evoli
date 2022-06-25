package component

import "github.com/evolidev/evoli/framework/view"

type Methods struct{}

func (c *Methods) Include(name string) string {
	return NewByName(name, nil).Render()
}

func SetupViewEngine(engine *view.Engine) {
	engine.AddRenderData("Component", &Methods{})
	engine.AddPlaceholder("@componentHeader", `
	<script src="https://cdn.tailwindcss.com"></script>
`)
	engine.AddPlaceholder("@componentFooter", `<!-- @componentFooter -->`)
}
