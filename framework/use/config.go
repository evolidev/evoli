package use

import (
	"github.com/evolidev/evoli/framework/config"
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
	conf := config.NewConfig(prefix)
	instances.Add(prefix, conf)

	instances.Add("test", config.NewConfig("test"))

	return conf
}

func ensureConfigPath() {
	if config.Directory() == "" {
		config.SetDirectory(BasePath("configs"))
	}
}
