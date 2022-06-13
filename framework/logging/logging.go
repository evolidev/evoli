package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/evolidev/evoli/framework/console/color"
	"github.com/mitchellh/go-homedir"
)

const logFormat = "%s"
const textColor = 242
const debugColor = 3
const successColor = 2
const errorColor = 1

type Logger struct {
	log *log.Logger
}

func NewLogger(c *Config) *Logger {
	var prefixColor = c.PrefixColor
	var w io.Writer = c.Stdout
	if w == nil {
		w = os.Stdout
	}

	return &Logger{
		log: log.New(w, color.Text(prefixColor, "["+c.Name+"] "), log.LstdFlags|log.Lmsgprefix),
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
			fmt.Sprintf("%s %s", color.Text(successColor, "Success"), color.Text(textColor, msg)),
			args...,
		),
	)
}

func (l *Logger) Error(msg interface{}, args ...interface{}) {
	l.log.Printf(
		fmt.Sprintf(
			fmt.Sprintf("%s %s", color.Text(errorColor, "Error"), color.Text(textColor, msg)),
			args...,
		),
	)
}

func (l *Logger) Debug(msg interface{}, args ...interface{}) {
	l.log.Printf(
		fmt.Sprintf(
			fmt.Sprintf("%s %s", color.Text(debugColor, "Debug"), color.Text(textColor, msg)),
			args...,
		),
	)
}

func (l *Logger) Print(msg interface{}, args ...interface{}) {
	l.log.Printf(fmt.Sprintf(logFormat, color.Text(textColor, msg)), args...)
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
