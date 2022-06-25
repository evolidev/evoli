package routes

import (
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

			return response.View("test").WithData(map[string]interface{}{
				"Name": persons[0].Name,
			})
		})

		personRouter.Post("/", func(form url.Values) *response.RedirectResponse {
			var p model.Person
			p.Name = form.Get("Name")
			use.DB().Create(&p)

			return response.Redirect("/person")
		})

		personRouter.Post("/:name", func(name string) *model.Person {
			person := model.Person{Name: name}
			use.DB().Create(&person)

			return &person
		})
	})

	web.Get("/component", func() *response.ViewResponse {
		return response.View("component")
	})
}
