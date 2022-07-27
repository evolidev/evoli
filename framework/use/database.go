package use

import (
	"fmt"
	"github.com/evolidev/evoli/framework/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/fs"
	"os"
	"strings"
)

var myDb *gorm.DB
var migrations Collection[int, func(db *gorm.DB)]
var migration *database.Migration

func Migration() *database.Migration {
	if nil == migration {
		migration = database.NewMigrations()
	}

	return migration
}

func DB() *gorm.DB {
	if nil != myDb {
		return myDb
	}

	config := Config("db.sqlite.path")

	dir := config.Value().(string)

	directories := strings.Split(dir, "/")
	directories = directories[:len(directories)-1]

	//todo get from storage
	tmp := os.DirFS(BasePath())

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

	db, _ := gorm.Open(sqlite.Open(dir), &gorm.Config{})

	myDb = db

	return myDb
}
