package use

import (
	"embed"
	"evoli.dev/framework/config"
	"strings"
)

func init() {
	instances = NewCollection[string, *config.Config]()
}

var instances *Collection[string, *config.Config]

func Config(path string) *config.Config {
	ensureConfigPath()

	items := strings.Split(path, ".")

	conf := addConfig(items[0])

	return conf.Get(strings.Join(items[1:], "."))
}

func addConfig(prefix string) *config.Config {
	if instances.Has(prefix) {
		return instances.Get(prefix)
	}

	store := Store("embed")
	if store != nil {
		storeFs := store.FS()
		if store, ok := storeFs.(embed.FS); ok {
			config.SetEmbed(store)
		}
	}

	conf := config.NewConfig(prefix)
	instances.Add(prefix, conf)

	return conf
}

func ensureConfigPath() {
	if config.Directory() == "" {
		config.SetDirectory(BasePath() + "configs")
	}
}
