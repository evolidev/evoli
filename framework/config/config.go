package config

import (
	"embed"
	"flag"
	"github.com/joho/godotenv"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var configDir string
var envRead = false
var embedFs embed.FS

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

func (c *Config) SetDefault(key string, value interface{}) *Config {
	c.instance.SetDefault(key, value)

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
	conf.SetFs(MyFS{afero.FromIOFS{FS: embedFs}})
	conf.AddConfigPath("/configs")
	conf.AutomaticEnv()
	conf.SetEnvKeyReplacer(getReplacer())
	err := conf.ReadInConfig()

	if err != nil {
		conf.AddConfigPath(configDir)
		conf.SetFs(afero.NewOsFs())
		err = conf.ReadInConfig()
		if err != nil {
			//fmt.Println(err)
		}
	}

	return &Config{instance: conf}
}

type MyFS struct {
	afero.FromIOFS
}

// Open will be overridden to strip the "/" away since viper will use an absolute path
func (receiver MyFS) Open(name string) (afero.File, error) {
	name = trimLeftChar(name)

	return receiver.FromIOFS.Open(name)
}

// Stat will be overridden to strip the "/" away since viper will use an absolute path
func (receiver MyFS) Stat(name string) (os.FileInfo, error) {
	name = trimLeftChar(name)

	return receiver.FromIOFS.Stat(name)
}

func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}

func SetDirectory(dir string) {
	configDir = dir
}

func SetEmbed(embed embed.FS) {
	embedFs = embed
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
