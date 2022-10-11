package logging

import "io"

type Config struct {
	EnableColors bool
	Name         string
	Stdout       io.Writer
	Path         string
	PrefixColor  int
}
