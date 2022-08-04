package filesystem

import "io/fs"

type Store interface {
	HasDir(dir string) bool
	FS() fs.FS
	Sub(dir string) fs.FS
}
