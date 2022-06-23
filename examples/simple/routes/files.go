package routes

import (
	"github.com/evolidev/evoli/examples/simple/resources"
	"github.com/evolidev/evoli/framework/router"
	"io/fs"
	"net/http"
)

func Files(files *router.Router) {
	sub, _ := fs.Sub(resources.Static, "static")
	files.ServeFiles("/static/*filepath", http.FS(sub))
}
