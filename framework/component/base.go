package component

import (
	"fmt"
	"github.com/evolidev/evoli/framework/response"
	"github.com/evolidev/evoli/framework/use"
)

type Base struct {
	File      string
	Component *use.Reflection
	Data      *use.Collection[string, interface{}]
}

func (b *Base) GetFilePath() string {
	if ok := b.Component.HasMethod("GetFilePath"); ok {
		output := b.Component.Method("GetFilePath").Call()
		return output.Interface().(string)
	}

	return fmt.Sprintf(
		"templates/%s",
		use.String(b.GetComponentName()).Kebab().Get(),
	)
}

func (b *Base) GetComponentName() string {
	return b.Component.Name()
}

func (b *Base) GetRawContent() string {
	path := b.GetFilePath()

	return string(response.View(path).AsBytes())
}

func (b *Base) Render() string {
	return b.GetRawContent()
}

func (b *Base) Set(data map[string]interface{}) {
	b.Component = b.Component.WithParams(data).Fill()
}

func (b *Base) Get(key string) interface{} {
	return b.Component.GetField(key)
}

func (b *Base) Call(method string, parameters interface{}) interface{} {
	output := b.Component.Method(method)

	result := output.WithParams(parameters)
	response := result.Call()
	return response.Interface()
}
