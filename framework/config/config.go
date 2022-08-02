package config

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var configDir string
var envRead = false

type Config struct {
	instance *viper.Viper
	key      string
}

func (c *Config) Get(key string) *Config {
	sub := c.instance.Sub(key)
	if nil == sub {
		sub = c.instance
	} else {
		key = ""
	}

	return &Config{instance: sub, key: key}
}

func (c *Config) Value() interface{} {
	if c.key == "" {
		return c
	}

	return c.instance.Get(c.key)
}

func (c *Config) Set(key string, value interface{}) *Config {
	c.instance.Set(key, value)

	return c
}

func NewConfig(prefix string) *Config {
	if !envRead {
		readEnv()
		envRead = true
	}
	conf := viper.New()
	conf.SetEnvPrefix(prefix)
	conf.SetConfigName(prefix)
	conf.AddConfigPath(configDir)
	conf.AutomaticEnv()
	conf.SetEnvKeyReplacer(getReplacer())
	conf.ReadInConfig()

	return &Config{instance: conf}
}

func SetDirectory(dir string) {
	configDir = dir
}

func Directory() string {
	return configDir
}

func getReplacer() *strings.Replacer {
	toReplace := make([]string, 2)
	toReplace = append(toReplace, ".")
	toReplace = append(toReplace, "_")

	return strings.NewReplacer(toReplace...)
}

func readEnv() {
	env := os.Getenv("APP_ENV")
	if "" == env {
		if flag.Lookup("test.v") == nil {
			env = "development"
		} else {
			env = "testing"
		}
	}

	godotenv.Load(".env." + env + ".local")
	if "test" != env {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load() // The Original .env
}
