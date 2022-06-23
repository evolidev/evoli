package logging

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/evolidev/evoli/framework/console/color"
	"github.com/mitchellh/go-homedir"
)

const logFormat = "%s"
const textColor = 245
const timeColor = 240
const debugColor = 3
const successColor = 2
const errorColor = 1
const logColor = 61

type Logger struct {
	log    *log.Logger
	config *Config
}

func NewLogger(c *Config) *Logger {
	var w = c.Stdout
	if w == nil {
		w = os.Stdout
	}

	return &Logger{
		log:    log.New(w, "", 0),
		config: c,
	}
}

func (l *Logger) getPrefix() string {
	var prefixColor = l.config.PrefixColor
	currentTime := time.Now()
	return fmt.Sprintf(
		"%s %s",
		color.Text(timeColor, currentTime.Format("2006-01-02 15:04:05")),
		color.Text(prefixColor, "["+l.config.Name+"]"),
	)
}

func (l *Logger) Log(msg interface{}, args ...interface{}) {
	l.log.Printf(
		fmt.Sprintf("%s %s %s", l.getPrefix(), color.Text(logColor, "Log"), color.Text(textColor, msg)),
		args...,
	)
}

func (l *Logger) Success(msg interface{}, args ...interface{}) {
	l.log.Printf(
		fmt.Sprintf("%s %s %s", l.getPrefix(), color.Text(successColor, "Success"), color.Text(textColor, msg)),
		args...,
	)
}

func (l *Logger) Error(msg interface{}, args ...interface{}) {
	l.log.Printf(
		fmt.Sprintf("%s %s %s", l.getPrefix(), color.Text(errorColor, "Error"), color.Text(textColor, msg)),
		args...,
	)
}

func (l *Logger) Debug(msg interface{}, args ...interface{}) {
	l.log.Printf(
		fmt.Sprintf("%s %s %s", l.getPrefix(), color.Text(debugColor, "Debug"), color.Text(textColor, msg)),
		args...,
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
