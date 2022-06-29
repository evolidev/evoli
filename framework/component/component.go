package component

import (
	"github.com/evolidev/evoli/framework/logging"
	"github.com/evolidev/evoli/framework/response"
	"github.com/evolidev/evoli/framework/router"
	"github.com/evolidev/evoli/framework/use"
)

const MOUNT = "Mount"
const UPDATE = "Update"

const ENDPOINT = "/internal/component"
const ASSET = "/vendor/evoli/static/component.js"

var components = make(map[string]Component)

type Component interface {
}

type Data struct {
	Name string `json:"name"`
}

type Request struct {
	Id         string         `json:"_id"`
	Component  string         `json:"component"`
	Method     string         `json:"method"`
	State      map[string]any `json:"state"`
	Action     string         `json:"action"`
	Parameters []any          `json:"parameters"`
}

type Response struct {
	Id        string         `json:"_id"`
	Component string         `json:"component"`
	State     map[string]any `json:"state"`
	Type      string         `json:"type"`
	Content   string         `json:"content"`
	Response  any            `json:"response"`
}

func New(componentStruct any, data map[string]any) *Base {
	collection := use.NewCollection[string, any]()
	collection.Set(data)

	component := use.Magic(componentStruct).ToPointer()

	return &Base{
		Component: component.WithParams(data).Fill(),
		Data:      collection,
	}
}

func Register(component Component) {
	name := use.GetInterfacedStructName(component)
	components[name] = component
}

func GetRegisterComponentsCount() int {
	return len(components)
}

func NewByNameWithData(name string, data string) *Base {
	componentObject, ok := components[name]

	if !ok {
		return nil
	}

	mappedData := use.JsonDecodeObject(data)
	component := New(componentObject, mappedData)

	return component
}

func NewByName(name string, data map[string]any) *Base {
	componentObject, ok := components[name]

	if !ok {
		return nil
	}

	component := New(componentObject, data)

	return component
}

func Handle(request *Request) *Response {
	component := NewByName(request.Component, request.State)

	if component == nil {
		return nil
	}

	var response any

	if request.Action == "click" {
		response = component.Call(request.Method, request.Parameters)
	}

	return &Response{
		Id:        request.Id,
		Component: request.Component,
		State:     component.GetState(),
		Response:  response,
	}
}

func RegisterRoutes(r *router.Router) {
	r.Post(ENDPOINT, handleRouterRequest)
	r.File("/vendor/evoli/static/component.js", "../../resources/component.js")
}

func handleRouterRequest(request *router.Request) any {
	r := request.Body()

	componentRequest := &Request{}
	use.JsonDecodeStruct(use.JsonEncode(r), componentRequest)

	if valid := validateRequest(componentRequest); !valid {
		logging.Error("Invalid request")
		return response.Json(map[string]any{"error": true}).WithCode(400)
	}

	res := Handle(componentRequest)

	if res == nil {
		return response.Json(map[string]any{"error": true}).WithCode(400)
	}

	return response.Json(res)
}

func validateRequest(request *Request) bool {
	if request.Component == "" {
		return false
	}

	if request.Action == "" {
		return false
	}

	return true
}
