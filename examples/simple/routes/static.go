package routes

import (
	"github.com/evolidev/evoli/framework/router"
)

func Files(r *router.Router) {
	r.File("/favicon.ico", "public/favicon.ico")
}
