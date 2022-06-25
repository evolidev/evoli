package use

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/fs"
	"os"
)

var myDb *gorm.DB

func DB() *gorm.DB {
	if nil != myDb {
		return myDb
	}

	tmp := os.DirFS("/")

	_, err := fs.ReadDir(tmp, "database")
	if err != nil {
		err = os.Mkdir("database", 0755)
		if err != nil {
			fmt.Println(err)
		}
	}
	db, _ := gorm.Open(sqlite.Open("database/test.db"), &gorm.Config{})

	myDb = db

	return myDb
}
