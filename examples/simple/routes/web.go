package routes

import (
	evoli "github.com/evolidev/evoli/framework"
	"github.com/evolidev/evoli/framework/response"
)

func Web(router *evoli.Router) {
	router.Get("/", func() string { return "hello" })
	router.Get("/test", func() *response.ViewResponse {
		return response.View("test").SetBasePath("examples/simple/resources/views/")
	})
}
