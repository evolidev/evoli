package component

import (
	"fmt"
	"github.com/evolidev/evoli/framework/use"
	"github.com/evolidev/evoli/framework/view"
	"html"
	"strings"
)

type Methods struct{}

func (c *Methods) Include(name string) string {
	component := NewByName(name, nil)
	rendered := component.Render()

	rendered = strings.ReplaceAll(
		rendered,
		"@scope",
		fmt.Sprintf(
			` data-cid="%s" data-scope="%s"`,
			component.GetCid(),
			html.EscapeString(use.JsonEncode(map[string]any{
				"state": component.GetState(),
				"name":  name,
			})),
		),
	)

	return rendered
}

func SetupViewEngine(engine *view.Engine) {
	engine.AddRenderData("Component", &Methods{})
	engine.AddPlaceholder("@componentHeader", `
	<script src="https://cdn.tailwindcss.com"></script>
`)
	engine.AddPlaceholder("@componentFooter", `
	<script src="https://unpkg.com/petite-vue"></script>
	<script src="/static/component.js"></script>
`)
}
