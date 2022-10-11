package test

import (
	"evoli.dev/framework/filesystem"
	"evoli.dev/framework/logging"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestLogging(t *testing.T) {
	t.Run("should log message", func(t *testing.T) {
		// create pipe
		r, w, _ := os.Pipe()
		defer r.Close()
		defer w.Close()

		logger := logging.NewLogger(&logging.Config{
			Stdout: w,
		})

		logger.Info("test")

		// read from pipe
		buf := make([]byte, 1024)
		n, _ := r.Read(buf)

		// check if message is logged
		assert.True(t, strings.Contains(string(buf[:n]), "test"), "message should be logged")
	})

	t.Run("should log message with prefix", func(t *testing.T) {
		// create pipe
		r, w, _ := os.Pipe()
		defer r.Close()
		defer w.Close()

		logger := logging.NewLogger(&logging.Config{
			Stdout: w,
			Name:   "prefix",
		})

		logger.Info("test")

		// read from pipe
		buf := make([]byte, 1024)
		n, _ := r.Read(buf)

		// check if message is logged
		assert.True(t, strings.Contains(string(buf[:n]), "prefix"), "message should be logged with prefix")
	})

	t.Run("should log to a file", func(t *testing.T) {
		logger := logging.NewLogger(&logging.Config{
			Stdout: os.Stdout,
			Path:   "test.log",
		})

		logger.Info("test loggger")

		// check if message is logged to file
		assert.FileExists(t, "test.log", "message should be logged to file")

		content := filesystem.Read("test.log")
		assert.True(t, strings.Contains(content, "test loggger"), "message should be logged to file")

		// remove file
		defer os.Remove("test.log")
	})
}
