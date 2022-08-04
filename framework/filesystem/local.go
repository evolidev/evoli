package filesystem

import (
	"github.com/evolidev/evoli/framework/config"
	"io/fs"
	"os"
)

type LocalFS struct {
	base string
}

func (f *LocalFS) HasDir(dir string) bool {
	matches, _ := fs.Glob(os.DirFS(f.base), dir)

	return len(matches) > 0
}

func NewLocalFS(config *config.Config) *LocalFS {
	return &LocalFS{base: config.Get("base_path").Value().(string)}
}
