package routes

import (
	"github.com/evolidev/evoli/framework/router"
)

func Folders(r *router.Router) {
	r.Static("/", "public")
	//r.File("/favicon.ico", "files/f")
}

func Files(r *router.Router) {
	r.File("/favicon.ico", "public/favicon.ico")
}
