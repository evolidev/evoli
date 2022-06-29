package component

import (
	"fmt"
	"github.com/evolidev/evoli/framework/response"
	"github.com/evolidev/evoli/framework/use"
	"github.com/matoous/go-nanoid/v2"
	"log"
)

type Base struct {
	Id        string
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
		"components/%s",
		use.String(b.GetComponentName()).Kebab().Get(),
	)
}

func (b *Base) GetComponentName() string {
	return b.Component.Name()
}

func (b *Base) GetRawContent() string {
	path := b.GetFilePath()

	return string(response.View(path).WithData(b.GetState()).AsBytes())
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
	if b.Component == nil {
		return nil
	}

	if !b.Component.HasMethod(method) {
		log.Printf("Component method not found: %s@%s", b.Component.Name(), method)
		return nil
	}

	output := b.Component.Method(method)

	result := output.WithParams(parameters)
	response := result.Call()

	if !response.IsValid() {
		return response
	}
	return response.Interface()
}

func (b *Base) GetState() map[string]interface{} {
	return b.Component.GetFields()
}

func (b *Base) Trigger(args ...any) {
	hook := args[0].(string)

	parameters := args[1:]

	if b.Component.HasMethod(hook) {
		b.Component.Method(hook).WithParams(parameters).Call()
	}
}

func (b *Base) GetData() map[string]any {
	return map[string]any{
		"state":     b.GetState(),
		"component": b.GetComponentName(),
		"_id":       b.GetCid(),
	}
}

func (b *Base) GetCid() string {
	if b.Id == "" {
		b.Id, _ = gonanoid.New()
	}

	return b.Id
}
