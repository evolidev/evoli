package routes

import (
	evoli "github.com/evolidev/evoli/framework"
)

func Web(router *evoli.Router) {
	router.Get("/", func() string { return "hello" })
	router.Get("/test", func() string { return "hello" })
}
