package test

import (
	"embed"
	"github.com/evolidev/evoli/framework/config"
	"github.com/evolidev/evoli/framework/filesystem"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

//go:embed resources
var tmp embed.FS

func TestEmbedFS(t *testing.T) {
	t.Run("has dir should return true if directory exists", func(t *testing.T) {
		embedFs := filesystem.NewEmbedFS(tmp)

		assert.True(t, embedFs.HasDir("resources/views"))
	})

	t.Run("has dir should return false if dir does not exists", func(t *testing.T) {
		embedFs := filesystem.NewEmbedFS(tmp)

		assert.False(t, embedFs.HasDir("not_exists"))
	})
}

func TestLocalFS(t *testing.T) {
	config.SetDirectory("./")
	cnf := config.NewConfig("storage")
	abs, _ := filepath.Abs("./")
	cnf.Set("local.base_path", abs)

	localStore := filesystem.NewLocalFS(cnf.Get("local"))

	assert.True(t, localStore.HasDir("resources/views"))
}
