package component

import (
	"fmt"
	"github.com/evolidev/evoli/framework/view"
)

type Methods struct{}

func (c *Methods) Include(name string, arg ...any) (string, error) {
	component := NewByName(name, nil)

	if component == nil {
		// TODO: Throw error
		return "", fmt.Errorf("Component %s not found", name)
	}

	args := append([]any{MOUNT}, arg...)
	component.Trigger(args...)

	return component.RenderParsed(), nil
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
