package filesystem

import (
	"embed"
	"io/fs"
)

type EmbedFS struct {
	fs embed.FS
}

func (f *EmbedFS) HasDir(dir string) bool {
	matches, _ := fs.Glob(f.fs, dir)

	return len(matches) > 0
}

func (f *EmbedFS) FS() fs.FS {
	return f.fs
}

func (f *EmbedFS) Sub(dir string) fs.FS {
	fis, _ := fs.Sub(f.fs, dir)

	return fis
}

func NewEmbedFS(fs embed.FS) *EmbedFS {
	return &EmbedFS{fs}
}
