package logging

import (
	"fmt"
	"github.com/evolidev/evoli/framework/console/color"
	"github.com/mitchellh/go-homedir"
	"io"
	"log"
	"os"
	"path"
)

const logFormat = "%s"

type Logger struct {
	log *log.Logger
}

func NewLogger(c *Config) *Logger {

	var w io.Writer = c.Stdout
	if w == nil {
		w = os.Stdout
	}

	return &Logger{
		log: log.New(w, color.Text(170, "["+c.Name+"] "), log.LstdFlags|log.Lmsgprefix),
	}
}

func (l *Logger) Log(color func(string, ...interface{}) string, prefix string, msg interface{}, args ...interface{}) {
	l.log.Print(
		fmt.Sprintf(color(prefix), msg, " ", args),
	)
}

func (l *Logger) Success(msg interface{}, args ...interface{}) {
	l.log.Printf(
		fmt.Sprintf(
			fmt.Sprintf("%s %s", color.Text(2, "Success"), color.Text(247, msg)),
			args...,
		),
	)
}

func (l *Logger) Error(msg interface{}, args ...interface{}) {
	l.log.Printf(
		fmt.Sprintf(
			fmt.Sprintf("%s %s", color.Text(1, "Error"), color.Text(247, msg)),
			args...,
		),
	)
}

func (l *Logger) Debug(msg interface{}, args ...interface{}) {
	l.log.Printf(
		fmt.Sprintf(
			fmt.Sprintf("%s %s", color.Text(3, "Debug"), color.Text(247, msg)),
			args...,
		),
	)
}

func (l *Logger) Print(msg interface{}, args ...interface{}) {
	l.log.Printf(fmt.Sprintf(logFormat, msg), args...)
}

var LogLocation = func() string {
	dir, _ := homedir.Dir()
	dir, _ = homedir.Expand(dir)
	dir = path.Join(dir, ".refresh")
	os.MkdirAll(dir, 0755)
	return dir
}

var ErrorLogPath = func() string {
	return path.Join(LogLocation(), "error.log")
}
