package test

import (
	"github.com/evolidev/evoli/framework/config"
	"github.com/evolidev/evoli/framework/use"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	config.SetDirectory("./configs")
	use.Embed(tmp)

	t.Run("config should get config value", func(t *testing.T) {
		conf := use.Config("storage")

		result := conf.Get("local.path")

		assert.Equal(t, "storage", result.Value())
	})

	t.Run("config should get config value directly", func(t *testing.T) {
		conf := use.Config("storage.local.path")

		result := conf.Value()

		assert.Equal(t, "storage", result)
	})

	t.Run("config should return sub config if key points to sub config", func(t *testing.T) {
		conf := use.Config("storage")

		result := conf.Get("local")

		if _, ok := result.Value().(*config.Config); !ok {
			t.Errorf("Not a config")
			t.Fail()
		}

		result = result.Value().(*config.Config).Get("path")
		assert.Equal(t, "storage", result.Value())
	})

	t.Run("config should be overridden by environment variable", func(t *testing.T) {
		os.Setenv("STORAGE_LOCAL_PATH", "test")
		conf := use.Config("storage.local.path")

		result := conf.Value()

		assert.Equal(t, "test", result)
	})
}

func TestSetConfig(t *testing.T) {
	config.SetDirectory("./configs")

	t.Run("config should set config value", func(t *testing.T) {
		conf := use.Config("storage")
		conf.Set("local.path", "test")

		result := conf.Get("local.path")

		assert.Equal(t, "test", result.Value())
	})
}
