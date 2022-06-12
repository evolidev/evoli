package test

import (
	"github.com/evolidev/evoli/framework/use"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestCall(t *testing.T) {
	t.Run("Call should call given function", func(t *testing.T) {
		m := use.Magic(func() string { return "test" })

		result := m.Call()

		assert.Exactly(t, "test", result.Interface().(string))
	})

	t.Run("Call should call structs pointer method", func(t *testing.T) {
		m := use.Magic(TestStructFirst{}).Method("TestPointer")

		result := m.Call()

		assert.Exactly(t, "success", result.Interface().(string))
	})

	t.Run("Call should call pointer structs pointer method", func(t *testing.T) {
		m := use.Magic(&TestStructFirst{}).Method("TestPointer")

		result := m.Call()

		assert.Exactly(t, "success", result.Interface().(string))
	})

	t.Run("Call should call structs value method", func(t *testing.T) {
		m := use.Magic(TestStructFirst{}).Method("TestValue")

		result := m.Call()

		assert.Exactly(t, "success", result.Interface().(string))
	})

	t.Run("Call should call direct structs pointer method", func(t *testing.T) {
		m := use.Magic((&TestStructFirst{}).TestPointer)

		result := m.Call()

		assert.Exactly(t, "success", result.Interface().(string))
	})

	t.Run("Call should call direct structs value method", func(t *testing.T) {
		m := use.Magic(TestStructFirst.TestValue)

		result := m.Call()

		assert.Exactly(t, "success", result.Interface().(string))
	})

	t.Run("Call should call direct structs pointer method and use its value", func(t *testing.T) {
		m := use.Magic((&TestStructFirst{"success"}).TestPointerWithValue)

		result := m.Call()

		assert.Exactly(t, "success", result.Interface().(string))
	})

	t.Run("Call should call given function with string parameter", func(t *testing.T) {
		m := use.Magic(func(test string) string { return "hello-" + test })
		params := make([]string, 0)
		params = append(params, "world")

		result := m.WithParams(params).Call()

		assert.Exactly(t, "hello-world", result.Interface().(string))
	})

	t.Run("Call should call given function with int parameter", func(t *testing.T) {
		m := use.Magic(func(test int) string { return "hello-" + strconv.Itoa(test) })
		params := make([]string, 0)
		params = append(params, "1")

		result := m.WithParams(params).Call()

		assert.Exactly(t, "hello-1", result.Interface().(string))
	})

	t.Run("Call should call given function with bool parameter", func(t *testing.T) {
		m := use.Magic(func(test bool) string { return "hello-" + strconv.FormatBool(test) })
		params := make([]string, 0)
		params = append(params, "1")

		result := m.WithParams(params).Call()

		assert.Exactly(t, "hello-true", result.Interface().(string))
	})

	t.Run("Call should call given function with mixed parameters", func(t *testing.T) {
		m := use.Magic(func(test string, myInt int) string { return "hello-" + test + "-" + strconv.Itoa(myInt) })
		params := make([]interface{}, 0)
		params = append(params, "world")
		params = append(params, 1)

		result := m.WithParams(params).Call()

		assert.Exactly(t, "hello-world-1", result.Interface().(string))
	})

	t.Run("Call should be able to change struct fields", func(t *testing.T) {
		myStruct := &TestStructFirst{"success"}
		m := use.Magic(myStruct)

		result := m.Method("SetTestPropWithPointer").Call()

		assert.Exactly(t, "updated", result.Interface().(string))
		assert.Exactly(t, "updated", myStruct.TestProp)
	})

	t.Run("Call should be able to change struct fields of value receiver if it is cast to pointer", func(t *testing.T) {
		myStruct := TestStructFirst{"success"}
		m := use.Magic(myStruct).ToPointer()

		result := m.Method("SetTestPropWithPointer").Call()

		assert.Exactly(t, "updated", result.Interface().(string))
		assert.Exactly(t, "updated", m.GetField("TestProp"))
	})
}

func TestNewPointer(t *testing.T) {
	t.Run("New pointer should return a pointer to struct", func(t *testing.T) {
		myStruct := TestStructFirst{"success"}
		m := use.Magic(myStruct)

		p := m.NewPointer().Interface()

		assert.Exactly(t, "", p.(*TestStructFirst).TestProp)
	})
}

func TestToPointer(t *testing.T) {
	t.Run("To pointer should do nothing if given struct is already a pointer", func(t *testing.T) {
		myStruct := &TestStructFirst{"success"}
		m := use.Magic(myStruct)

		result := m.ToPointer()

		assert.Exactly(t, "success", result.Value().(*TestStructFirst).TestProp)
	})

	t.Run("To pointer should convert current struct to pointer", func(t *testing.T) {
		myStruct := TestStructFirst{"success"}
		m := use.Magic(myStruct).ToPointer()

		//result := m.Method("SetTestPropWithPointer").Call()

		assert.Exactly(t, "success", m.Value().(*TestStructFirst).TestProp)
	})
}

func TestFill(t *testing.T) {
	t.Run("Fill should fill data to struct", func(t *testing.T) {
		params := make(map[string]interface{})
		params["TestProp"] = "test"
		m := use.Magic(&TestStructFirst{})

		result := m.WithParams(params).Fill().Value().(*TestStructFirst)

		assert.Exactly(t, "test", result.TestProp)
	})
}

func test(test1 string) {

}

type TestStructFirst struct {
	TestProp string
}

func (receiver TestStructFirst) TestPointerWithValue() string {
	return receiver.TestProp
}

func (receiver *TestStructFirst) SetTestPropWithPointer() string {
	receiver.TestProp = "updated"

	return receiver.TestProp
}

func (receiver *TestStructFirst) TestPointer() string {
	return "success"
}

func (receiver TestStructFirst) TestValue() string {
	return "success"
}

type TestStructSecond struct {
	f interface{}
}

func New(t interface{}) {

}
