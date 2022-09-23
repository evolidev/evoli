package router

import (
	"fmt"
	"github.com/evolidev/evoli/framework/use"
	"github.com/evolidev/evoli/framework/validation"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

type Request struct {
	*http.Request
	p params
}

func (r *Request) Form() url.Values {
	return r.p.form
}

func (r *Request) RouteParams() httprouter.Params {
	return r.p.route
}

func (r *Request) Query() url.Values {
	return r.p.query
}

func (r *Request) Body() any {
	return r.p.body
}

func (r *Request) Get(param string) interface{} {
	if r.p.route.ByName(param) != "" {
		return r.p.route.ByName(param)
	}
	if r.p.query.Has(param) {
		return r.p.query.Get(param)
	}
	if r.p.form.Has(param) {
		return r.p.form.Get(param)
	}

	body := reflect.ValueOf(r.p.body)
	if body.Kind() == reflect.Map {
		result := body.MapIndex(reflect.ValueOf(param))
		if result.IsValid() {
			return result.Interface()
		}
	}
	if body.Kind() == reflect.Struct {
		if body.FieldByName(param).IsValid() {
			return body.FieldByName(param).Interface()
		}
	}

	return nil
}

func (r *Request) Params() *use.Collection[string, interface{}] {
	p := use.NewCollection[string, interface{}]()
	p.Add("Request", r)
	p.Add("HttpRequest", r.Request)
	p.Add("Form", r.p.form)
	p.Add("Query", r.p.query)
	p.Add("Route", r.p.route)

	return p
}

func (r *Request) ValidateStruct(s interface{}) validation.Validator {
	validator := validation.NewValidator()

	return validator.ValidateStruct(s)
}

func (r *Request) Validate(rules map[string]interface{}) validation.Validator {
	validator := validation.NewValidator()
	form := r.Params().Get("Form").(url.Values)
	data := make(map[string]interface{})

	for key, _ := range form {
		data[key] = form.Get(key)
	}

	return validator.Validate(data, rules)
}

func NewRequest(r *http.Request) *Request {
	err := r.ParseForm()
	if err != nil {
		// todo log into our logger
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//lm.logger.Error(err)
	}

	myForm := r.PostForm
	if len(myForm) == 0 {
		myForm = r.Form
	}

	p := params{
		route: httprouter.ParamsFromContext(r.Context()),
		form:  myForm,
		query: r.URL.Query(),
		body:  use.JsonDecode(string(body)),
	}

	return &Request{Request: r, p: p}
}

type params struct {
	route httprouter.Params
	query url.Values
	form  url.Values
	body  interface{}
}
