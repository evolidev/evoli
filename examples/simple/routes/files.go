package routes

import (
	"github.com/evolidev/evoli/framework/router"
)

func Files(r *router.Router) {
	r.Static("/", "public")
	//r.File("/favicon.ico", "files/f")
}
