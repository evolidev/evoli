package component

import (
	"fmt"
	"github.com/evolidev/evoli/framework/view"
)

type Methods struct{}

func (c *Methods) Include(name string, arg ...any) string {
	component := NewByName(name, nil)

	args := append([]any{MOUNT}, arg...)
	component.Trigger(args...)

	return component.RenderParsed()
}

func SetupViewEngine(engine *view.Engine) {
	engine.AddRenderData("Component", &Methods{})
	engine.AddPlaceholder("@componentHeader", `
	<script src="https://cdn.tailwindcss.com"></script>
`)
	engine.AddPlaceholder("@componentFooter", fmt.Sprintf(`
	<script src="https://unpkg.com/evoli-petite-vue@0.0.3"></script>
	<script src="%s"></script>
	`, ASSET))
}
