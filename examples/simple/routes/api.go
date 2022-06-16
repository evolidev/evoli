package routes

import (
	evoli "github.com/evolidev/evoli/framework"
	"github.com/evolidev/evoli/framework/response"
)

func Api(router *evoli.Router) {
	router.Get("/", func() *response.JsonResponse { return response.Json([]string{"hi"}) })
	router.Get("/test", func() struct{ Test string } { return struct{ Test string }{"test"} })
}
