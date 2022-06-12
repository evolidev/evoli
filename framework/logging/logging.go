package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
)

const logFormat = "%s"

type Logger struct {
	log *log.Logger
}

func NewLogger(c *Config) *Logger {
	color.NoColor = !c.EnableColors

	if runtime.GOOS == "windows" {
		color.NoColor = true
	}

	var w io.Writer = c.Stdout
	if w == nil {
		w = os.Stdout
	}

	return &Logger{
		log: log.New(w, fmt.Sprintf("[%s] ", c.Name), log.LstdFlags|log.Lmsgprefix),
	}
}

func (l *Logger) Log(color func(string, ...interface{}) string, prefix string, msg interface{}, args ...interface{}) {
	l.log.Print(
		fmt.Sprintf(color(prefix), msg, " ", args),
	)
}

func (l *Logger) Success(msg interface{}, args ...interface{}) {
	l.Log(color.GreenString, "Success", msg, args...)
}

func (l *Logger) Error(msg interface{}, args ...interface{}) {
	l.Log(color.RedString, "Error", msg, args...)
}

func (l *Logger) Debug(msg interface{}, args ...interface{}) {
	l.Log(color.YellowString, "Debug", msg, args...)
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
