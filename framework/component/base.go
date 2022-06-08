package component

import (
	"fmt"
	"github.com/evolidev/evoli/framework/filesystem"
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
		"templates/%s.html",
		use.String(b.GetComponentName()).Kebab().Get(),
	)
}

func (b *Base) GetComponentName() string {
	return b.Component.Name()
}

func (b *Base) GetRawContent() string {
	path := b.GetFilePath()

	return filesystem.Read(path)
}

func (b *Base) Render() string {
	return b.GetRawContent()
}

func (b *Base) Set(data map[string]interface{}) {
	b.Data.Set(data)
}

func (b *Base) Get(key string) interface{} {
	return b.Data.Get(key)
}

func (b *Base) Call(method string, parameters interface{}) interface{} {
	output := b.Component.Method(method)

	result := output.WithParams(parameters)
	response := result.Call()
	return response.Interface()
	//return b.Component.Method(method).Call().Interface()
}