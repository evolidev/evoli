package component

import (
	"fmt"
	"github.com/evolidev/evoli/framework/filesystem"
	"github.com/evolidev/evoli/framework/use"
	"log"
	"reflect"
)

type Base struct {
	File      string
	Component *Component
	Data      *use.Collection[string, interface{}]
}

func (b *Base) GetFilePath() string {
	if ok, method := use.HasMethod(b.Component, "GetFilePath"); ok {
		output := method.Call([]reflect.Value{})
		return output[0].String()
	}

	return fmt.Sprintf(
		"templates/%s.html",
		use.String(b.GetComponentName()).Kebab().Get(),
	)
}

func (b *Base) GetComponentName() string {
	return use.GetInterfacedStructName(*b.Component)
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

func (b *Base) Call(method string, parameters []interface{}) interface{} {
	use.HasMethod(b.Component, method)
	use.HasMethod(*b.Component, method)
	use.HasMethod(&b.Component, method)

	if ok, method := use.HasMethod(&b.Component, method); ok {
		use.P("Method found")
		output := method.Call([]reflect.Value{})
		return output[0].String()
	}

	log.Println("Failedeeeee", method, parameters)
	use.D(b.Component)

	return nil
}
