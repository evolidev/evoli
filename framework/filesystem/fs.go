package filesystem

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"net/http"
	"time"
)

type FS struct {
	fs fs.FS
}

func (f *FS) HasDir(dir string) bool {
	matches, _ := fs.Glob(f.fs, dir)

	return len(matches) > 0
}

func (f *FS) HttpFS() http.FileSystem {
	return http.FS(f.fs)
}

func (f *FS) FS() fs.FS {
	return f.fs
}

func (f *FS) ServeContent(writer http.ResponseWriter, request *http.Request, file string) {
	data, _ := ioutil.ReadFile(file)
	http.ServeContent(writer, request, file, time.Now(), bytes.NewReader(data))
}

func (f *FS) Sub(path string) *FS {
	sub, _ := fs.Sub(f.fs, path)
	return NewFS(sub)
}

func NewFS(fs fs.FS) *FS {
	return &FS{fs}
}
