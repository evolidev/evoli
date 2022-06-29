package component

import (
	"fmt"
	"github.com/evolidev/evoli/framework/use"
	"github.com/evolidev/evoli/framework/view"
	"html"
	"strings"
)

type Methods struct{}

func (c *Methods) Include(name string, arg ...any) string {
	component := NewByName(name, nil)

	if len(arg) > 0 {
		args := append([]any{MOUNT}, arg...)

		component.Trigger(args...)
	}

	rendered := component.Render()

	rendered = strings.ReplaceAll(
		rendered,
		"@scope",
		fmt.Sprintf(
			`v-scope="mount(%s)"`,
			html.EscapeString(use.JsonEncode(component.GetData())),
		),
	)

	return rendered
}

func SetupViewEngine(engine *view.Engine) {
	engine.AddRenderData("Component", &Methods{})
	engine.AddPlaceholder("@componentHeader", `
	<script src="https://cdn.tailwindcss.com"></script>
`)
	engine.AddPlaceholder("@componentFooter", fmt.Sprintf(`
	<script src="https://unpkg.com/evoli-petite-vue"></script>
	<script src="%s"></script>
	`, ASSET))
}
