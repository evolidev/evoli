package filesystem

import (
	"io/fs"
	"net/http"
)

type Store interface {
	HasDir(dir string) bool
	HttpFS() http.FileSystem
	FS() fs.FS
	ServeContent(writer http.ResponseWriter, request *http.Request, file string)
	Sub(path string) *FS
}
