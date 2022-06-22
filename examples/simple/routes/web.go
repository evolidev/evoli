package routes

import (
	"fmt"
	"github.com/evolidev/evoli/examples/simple/model"
	"github.com/evolidev/evoli/framework/response"
	"github.com/evolidev/evoli/framework/router"
	"github.com/evolidev/evoli/framework/use"
	"net/url"
)

func Web(web *router.Router) {
	web.Get("/", func() string { return "hello" })

	web.Prefix("/person").Group(func(personRouter *router.Router) {
		personRouter.Get("/", func() *response.ViewResponse {
			var persons []model.Person
			use.DB().Find(&persons)
			fmt.Println(persons)

			return response.View("test").WithData(persons)
		})

		personRouter.Post("/", func(form url.Values) *response.ViewResponse {
			var p model.Person
			p.Name = form.Get("Name")
			use.DB().Create(&p)

			var persons []model.Person
			use.DB().Find(&persons)
			fmt.Println(persons)

			return response.View("test").WithData(persons)
		})

		personRouter.Post("/:name", func(name string) *model.Person {
			person := model.Person{Name: name}
			use.DB().Create(&person)

			return &person
		})
	})
}
