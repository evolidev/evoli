package database

import (
	"fmt"
	"github.com/evolidev/evoli/framework/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/fs"
	"os"
	"strings"
)

type Sqlite struct {
	config *config.Config
}

func (s *Sqlite) Open() gorm.Dialector {
	s.config.SetDefault("path", "database/db.sqlite")
	dir := s.config.Get("path").Value().(string)

	directories := strings.Split(dir, "/")
	directories = directories[:len(directories)-1]

	tmp := os.DirFS(s.config.Get("base").Value().(string))

	for _, d := range directories {
		_, err := fs.ReadDir(tmp, d)
		if err != nil {
			fmt.Println(err)
			err = os.Mkdir(d, 0755)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	return sqlite.Open(dir)
}
