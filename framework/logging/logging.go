package logging

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/evolidev/evoli/framework/console/color"
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
	if c == nil {
		c = &Config{}
	}

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

func (l *Logger) Fatal(msg interface{}, args ...interface{}) {
	l.log.Printf(fmt.Sprintf(logFormat, color.Text(textColor, msg)), args...)
	os.Exit(1)
}

var appLogger = NewLogger(nil)

func GetAppLogger() *Logger {
	return appLogger
}

func SetAppLogger(l *Logger) {
	appLogger = l
}

func Debug(msg interface{}, args ...interface{}) {
	appLogger.Debug(msg, args...)
}

func Error(msg interface{}, args ...interface{}) {
	appLogger.Error(msg, args...)
}

func Fatal(msg interface{}, args ...interface{}) {
	appLogger.Fatal(msg, args...)
}

func Log(msg interface{}, args ...interface{}) {
	appLogger.Log(msg, args...)
}

func Success(msg interface{}, args ...interface{}) {
	appLogger.Success(msg, args...)
}
