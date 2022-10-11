package logging

import (
	"fmt"
	"github.com/evolidev/evoli/framework/filesystem"
	"io"
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
	log      *log.Logger
	plainLog *log.Logger
	config   *Config
}

var Verbose = 0

func NewLogger(c *Config) *Logger {
	if c == nil {
		c = &Config{
			Name: "app",
		}
	}

	var colorfulWriters []io.Writer
	var plainWriters []io.Writer

	if c.Stdout != nil {
		colorfulWriters = append(colorfulWriters, c.Stdout)
	} else {
		colorfulWriters = append(colorfulWriters, os.Stdout)
	}

	if c.Location != "" {
		if !filesystem.Exists(c.Location) {
			filesystem.Write(c.Location, "")
		}

		f, err := os.OpenFile(c.Location, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}

		plainWriters = append(plainWriters, f)
	}

	return &Logger{
		log:      log.New(io.MultiWriter(colorfulWriters...), "", 0),
		plainLog: log.New(io.MultiWriter(plainWriters...), "", 0),
		config:   c,
	}
}

func NewLoggerByName(name string, colorCode int) *Logger {
	return NewLogger(&Config{
		Name:        name,
		PrefixColor: colorCode,
	})
}

func (l *Logger) getPrefix() string {
	var prefixColor = l.config.PrefixColor
	currentTime := time.Now()
	prefix := ""

	if l.config.Name != "" {
		prefix = color.Text(prefixColor, "["+l.config.Name+"]") + " "
	}

	return fmt.Sprintf(
		"%s %s",
		color.Text(timeColor, currentTime.Format("2006-01-02 15:04:05")),
		prefix,
	)
}

func (l *Logger) getPlainPrefix() string {
	currentTime := time.Now()
	prefix := ""

	if l.config.Name != "" {
		prefix = "[" + l.config.Name + "]" + " "
	}

	return fmt.Sprintf(
		"%s %s",
		currentTime.Format("2006-01-02 15:04:05"),
		prefix,
	)
}

func (l *Logger) Log(msg interface{}, args ...interface{}) {
	l.log.Printf(
		fmt.Sprintf("%s%s %s", l.getPrefix(), color.Text(logColor, "Log"), color.Text(textColor, msg)),
		args...,
	)

	l.plainLog.Printf(
		fmt.Sprintf("%s%s %s", l.getPlainPrefix(), "Log", msg),
		args...,
	)
}

func (l *Logger) Info(msg interface{}, args ...interface{}) {
	l.log.Printf(
		fmt.Sprintf("%s%s %s", l.getPrefix(), color.Text(logColor, "Info"), color.Text(textColor, msg)),
		args...,
	)

	l.plainLog.Printf(
		fmt.Sprintf("%s%s %s", l.getPlainPrefix(), "Info", msg),
		args...,
	)
}

func (l *Logger) Success(msg interface{}, args ...interface{}) {
	l.log.Printf(
		fmt.Sprintf("%s%s %s", l.getPrefix(), color.Text(successColor, "Success"), color.Text(textColor, msg)),
		args...,
	)

	l.plainLog.Printf(
		fmt.Sprintf("%s%s %s", l.getPlainPrefix(), "Success", msg),
		args...,
	)
}

func (l *Logger) Error(msg interface{}, args ...interface{}) {
	l.log.Printf(
		fmt.Sprintf("%s%s %s", l.getPrefix(), color.Text(errorColor, "Error"), color.Text(textColor, msg)),
		args...,
	)

	l.plainLog.Printf(
		fmt.Sprintf("%s%s %s", l.getPlainPrefix(), "Error", msg),
		args...,
	)
}

func (l *Logger) Debug(msg interface{}, args ...interface{}) {
	l.log.Printf(
		fmt.Sprintf("%s%s %s", l.getPrefix(), color.Text(debugColor, "Debug"), color.Text(textColor, msg)),
		args...,
	)

	l.plainLog.Printf(
		fmt.Sprintf("%s%s %s", l.getPlainPrefix(), "Debug", msg),
		args...,
	)
}

func (l *Logger) Print(msg interface{}, args ...interface{}) {
	l.log.Printf(fmt.Sprintf(logFormat, color.Text(textColor, msg)), args...)

	l.plainLog.Printf(fmt.Sprintf(logFormat, msg), args...)
}

func (l *Logger) Fatal(msg interface{}, args ...interface{}) {
	l.log.Printf(fmt.Sprintf(logFormat, color.Text(textColor, msg)), args...)

	l.plainLog.Printf(fmt.Sprintf(logFormat, msg), args...)

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

func Info(msg interface{}, args ...interface{}) {
	appLogger.Info(msg, args...)
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
