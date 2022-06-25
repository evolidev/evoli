package reload

import (
	"context"
	"errors"
	"log"
	"os"
)

func Init() {
	log.Println("hello")
}

// ErrConfigNotExist is returned when a configuration file cannot be found.
var ErrConfigNotExist = errors.New("no config file was found")
var debug = true

func RunBackground(config *Configuration) error {
	r := NewWithContext(config, context.Background())
	return r.Start()
}

func Run(cfgFile string) error {
	ctx := context.Background()
	return RunWithContext(cfgFile, ctx)
}

func RunWithContext(cfgFile string, ctx context.Context) error {
	c := &Configuration{}

	if err := loadConfig(c, cfgFile); err != nil {
		if err != ErrConfigNotExist {
			return err
		}

		log.Println("No configuration loaded, proceeding with defaults")
	}

	if len(c.Path) > 0 {
		log.Printf("Configuration loaded from %s\n", c.Path)
	}

	if debug {
		c.Debug = true
	}

	r := NewWithContext(c, ctx)
	return r.Start()
}

func loadConfig(c *Configuration, path string) error {
	if len(path) > 0 {
		return c.Load(path)
	}

	for _, f := range [4]string{
		".refresh.yml",
		".refresh.yaml",
		"refresh.yml",
		"refresh.yaml",
	} {
		err := c.Load(f)
		if err != nil && os.IsNotExist(err) {
			continue
		}
		return err
	}

	return ErrConfigNotExist
}

//var cfgFile string

//func createConfig() {
//	c := Configuration{
//		AppRoot:            ".",
//		IgnoredFolders:     []string{"vendor", "log", "logs", "tmp", "node_modules", "bin", "templates"},
//		IncludedExtensions: []string{".go"},
//		BuildTargetPath:    "",
//		BuildPath:          os.TempDir(),
//		BuildDelay:         200,
//		BinaryName:         "refresh-build",
//		CommandFlags:       []string{},
//		CommandEnv:         []string{},
//		EnableColors:       true,
//	}
//
//	if cfgFile == "" {
//		cfgFile = "refresh.yml"
//	}
//
//	_, err := os.Stat(cfgFile)
//	if !os.IsNotExist(err) {
//		fmt.Errorf("config file %q already exists, skipping init", cfgFile)
//	}
//
//	c.Dump(cfgFile)
//}
