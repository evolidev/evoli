package routes

import (
	"github.com/evolidev/evoli/framework/response"
	"github.com/evolidev/evoli/framework/router"
)

func Web(router *router.Router) {
	router.Get("/", func() string { return "hello" })
	router.Get("/test", func() *response.ViewResponse {
		return response.View("test")
	})
}
