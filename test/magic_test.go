package test

import (
	"github.com/evolidev/evoli/framework/use"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCall(t *testing.T) {
	t.Run("Call should call given function", func(t *testing.T) {
		m := use.Magic(func() string { return "test" })

		result := m.Call()

		assert.Exactly(t, "test", result.Interface().(string))
	})

	t.Run("Call should call given function with parameters", func(t *testing.T) {
		m := use.Magic(func(test string) string { return "hello-" + test })
		params := make([]string, 0)
		params = append(params, "world")

		result := m.WithParams(params).Call()

		assert.Exactly(t, "hello-world", result.Interface().(string))
	})

	t.Run("asdfasdf", func(t *testing.T) {
		tmp := &TestStructSecond{&TestStructFirst{}}
		m := use.Magic(tmp.f)

		result := m.Method("Test").Call()

		assert.Exactly(t, "success", result.Interface().(string))
	})

}

func test(test1 string) {

}

type TestStructFirst struct {
}

func (receiver *TestStructFirst) Test() string {
	return "success"
}

type TestStructSecond struct {
	f interface{}
}

func New(t interface{}) {

}
