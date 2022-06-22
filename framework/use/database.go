package use

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var myDb *gorm.DB

func DB() *gorm.DB {
	if nil != myDb {
		return myDb
	}

	db, _ := gorm.Open(sqlite.Open("database/test.db"), &gorm.Config{})

	myDb = db

	return myDb
}
