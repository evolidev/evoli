package logging

import "io"

type Config struct {
	EnableColors bool
	Name         string
	Stdout       io.Writer
	Location     string
	PrefixColor  int
}
