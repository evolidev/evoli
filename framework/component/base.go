package component

import (
	"fmt"
	"github.com/evolidev/evoli/framework/filesystem"
	"github.com/evolidev/evoli/framework/use"
	"reflect"
)

type Base struct {
	File               string
	Component          Component
	ComponentInterface interface{}
	Data               *use.Collection[string, interface{}]
}

func (b *Base) GetComponentInterface() interface{} {
	if b.ComponentInterface != nil {
		return b.ComponentInterface
	}

	componentInterface := reflect.New(reflect.TypeOf(b.Component)).Interface()

	b.ComponentInterface = componentInterface

	return b.ComponentInterface
}

func (b *Base) GetFilePath() string {
	if ok, method := use.HasMethod(b.Component, "GetFilePath"); ok {
		output := method.Call([]reflect.Value{})
		return output[0].String()
	}

	return fmt.Sprintf("templates/%s.html", use.String(b.GetComponentName()).Kebab().Get())
}

func (b *Base) GetComponentName() string {
	return use.GetInterfacedStructName(b.Component)
}

func (b *Base) GetRawContent() string {
	path := b.GetFilePath()
	return filesystem.Read(path)
}

func (b *Base) Render() string {
	content := b.GetRawContent()
	return content
}

func (b *Base) Set(data map[string]interface{}) {
	b.Data.Set(data)
}

func (b *Base) Get(key string) interface{} {
	return b.Data.Get(key)
}
